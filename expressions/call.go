package expressions

import (
	"fmt"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type CallExpression struct {
	Pos  lexer.Position
	Name string
	Args []Expression
}

func (x CallExpression) Position() lexer.Position {
	return x.Pos
}

func (x CallExpression) Typecheck(inputShape Shape, params []*params.Param) (Shape, states.Action, *states.IDStack, error) {
	// go down the function stack and find the function invoked by this
	// call
	stack := inputShape.Stack
	for {
		// reached bottom of stack without finding a matching funcer
		if stack == nil {
			return Shape{}, nil, nil, errors.TypeError(
				errors.Code(errors.NoSuchFunction),
				errors.Pos(x.Pos),
				errors.InputType(inputShape.Type),
				errors.Name(x.Name),
				errors.NumParams(len(x.Args)+len(params)),
			)
		}
		// try the funcer on top of the stack
		funcer := stack.Head
		funOutputShape, funAction, ids, ok, err := funcer(inputShape, x, params)
		if err != nil {
			return Shape{}, nil, nil, err
		}
		if !ok {
			stack = stack.Tail
			continue
		}
		// make call lazy
		action := func(inputState states.State, args []states.Action) *states.Thunk {
			return &states.Thunk{
				Func: func() *states.Thunk {
					return funAction(inputState, args)
				},
				Stack:     inputState.Stack,
				TypeStack: inputState.TypeStack,
			}
		}
		// return
		return funOutputShape, action, ids, nil
	}
}

type Shape struct {
	Type  types.Type
	Stack *FuncerStack
}

type FuncerStack struct {
	Head Funcer
	Tail *FuncerStack
}

func (s *FuncerStack) Push(funcer Funcer) *FuncerStack {
	return &FuncerStack{funcer, s}
}

func (s *FuncerStack) PushAll(funcers []Funcer) *FuncerStack {
	for i := len(funcers) - 1; i >= 0; i-- {
		s = s.Push(funcers[i])
	}
	return s
}

func (s *FuncerStack) String() string {
	slice := make([]Funcer, 0)
	stack := s
	for stack != nil {
		slice = append(slice, stack.Head)
		stack = stack.Tail
	}
	return fmt.Sprintf("%v", slice)
}

func instantiate(pars []*params.Param, bindings map[string]types.Type) []*params.Param {
	result := make([]*(params.Param), len(pars))
	for i, par := range pars {
		result[i] = par.Instantiate(bindings)
	}
	return result
}

type Funcer func(gotInputShape Shape, gotCall CallExpression, gotParams []*params.Param) (outputShape Shape, action states.Action, ids *states.IDStack, ok bool, err error)

type SimpleKernel func(inputValue states.Value, argValues []states.Value) (states.Value, error)

func SimpleFuncer(wantInputType types.Type, wantName string, argTypes []types.Type, outputType types.Type, simpleKernel SimpleKernel) Funcer {
	// make parameters from argument types
	pars := make([]*params.Param, len(argTypes))
	for i, argType := range argTypes {
		pars[i] = &params.Param{
			InputType:  types.Any{},
			Params:     nil,
			OutputType: argType,
		}
	}
	// make regular kernel from simple kernel
	regularKernel := func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		argValues := make([]states.Value, len(argTypes))
		for i, arg := range args {
			res := arg(inputState.Clear(), nil).Eval()
			if res.Error != nil {
				return states.ThunkFromError(res.Error)
			}
			argValues[i] = res.Value
		}
		value, err := simpleKernel(inputState.Value, argValues)
		if err != nil {
			return states.ThunkFromError(err)

		}
		return states.ThunkFromState(states.State{
			Value:     value,
			Stack:     inputState.Stack,
			TypeStack: inputState.TypeStack,
		})

	}
	// return
	return RegularFuncer(wantInputType, wantName, pars, outputType, regularKernel, nil)
}

func VariableFuncer(id interface{}, name string, varType types.Type) Funcer {
	kernel := func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		stack := inputState.Stack
		for stack != nil {
			if stack.Head.ID == id {
				res := stack.Head.Action(states.InitialState, nil).Eval()
				if res.Error != nil {
					return states.ThunkFromError(res.Error)
				}
				return states.ThunkFromState(states.State{
					Value:     res.Value,
					Stack:     inputState.Stack,
					TypeStack: inputState.TypeStack,
				})

			}
			stack = stack.Tail
		}
		panic(fmt.Sprintf("variable %s not found", name))
	}
	return RegularFuncer(types.Any{}, name, nil, varType, kernel, &states.IDStack{
		Head: id,
	})
}

type RegularKernel func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk

func RegularFuncer(wantInputType types.Type, wantName string, pars []*params.Param, outputType types.Type, kernel RegularKernel, ids *states.IDStack) Funcer {
	return func(gotInputShape Shape, gotCall CallExpression, gotParams []*params.Param) (Shape, states.Action, *states.IDStack, bool, error) {
		// match number of parameters
		if len(gotCall.Args)+len(gotParams) != len(pars) {
			return Shape{}, nil, nil, false, nil
		}
		// match name
		if gotCall.Name != wantName {
			return Shape{}, nil, nil, false, nil
		}
		// match input type
		bindings := make(map[string]types.Type)
		if !wantInputType.Bind(gotInputShape.Type, bindings) {
			return Shape{}, nil, nil, false, nil
		}
		if !wantInputType.Instantiate(bindings).Subsumes(gotInputShape.Type) {
			return Shape{}, nil, nil, false, nil
		}
		// typecheck and set parameters filled by this call
		funAction := func(inputState states.State, args []states.Action) *states.Thunk {
			return kernel(inputState, args, bindings, gotCall.Position())
		}
		argActions := make([]states.Action, len(gotCall.Args))
		argIDss := make([]*states.IDStack, len(gotCall.Args))
		for i := range gotCall.Args {
			argInputShape := Shape{
				Type: pars[i].InputType.Instantiate(bindings),
				// TODO what if we don't have the bindings yet?
				// TODO what if an incompatible bound is declared?
				Stack: gotInputShape.Stack,
			}
			argOutputShape, argAction, argIDs, err := gotCall.Args[i].Typecheck(argInputShape, instantiate(pars[i].Params, bindings))
			if err != nil {
				return Shape{}, nil, nil, false, err
			}
			if !pars[i].OutputType.Bind(argOutputShape.Type, bindings) {
				return Shape{}, nil, nil, false, errors.TypeError(
					errors.Code(errors.ArgHasWrongOutputType),
					errors.Pos(gotCall.Pos),
					errors.ArgNum(i+1),
					errors.WantType(pars[i].OutputType.Instantiate(bindings)),
					errors.GotType(argOutputShape.Type),
				)
			}
			if !pars[i].OutputType.Instantiate(bindings).Subsumes(argOutputShape.Type) {
				return Shape{}, nil, nil, false, errors.TypeError(
					errors.Code(errors.ArgHasWrongOutputType),
					errors.Pos(gotCall.Pos),
					errors.ArgNum(i+1),
					errors.WantType(pars[i].OutputType.Instantiate(bindings)),
					errors.GotType(argOutputShape.Type),
				)
			}
			argActions[i] = argAction
			argIDss[i] = argIDs
			ids = ids.AddAll(argIDs)
		}
		// pass input variable stack to arguments
		funAction2 := func(inputState states.State, args []states.Action) *states.Thunk {
			args2 := make([]states.Action, len(argActions)+len(args))
			for i := range argActions {
				prunedStack := inputState.Stack.Keep(argIDss[i]) // pass only what is needed to avoid memory leak
				i := i
				args2[i] = func(argInputState states.State, argArgs []states.Action) *states.Thunk {
					argInputState.Stack = prunedStack
					return argActions[i](argInputState, argArgs)
				}
			}
			for i := 0; i < len(args); i++ {
				args2[len(argActions)+i] = args[i]
			}
			return funAction(inputState, args2)
		}
		// typecheck parameters not filled by the call
		for i, gotParam := range gotParams {
			wantParam := pars[len(gotCall.Args)+i].Instantiate(bindings)
			if !gotParam.Subsumes(wantParam) {
				return Shape{}, nil, nil, false, errors.TypeError(
					errors.Code(errors.ParamDoesNotMatch),
					errors.Pos(gotCall.Pos),
					errors.ParamNum(i+1),
					errors.WantParam(gotParam),
					errors.GotParam(wantParam),
				)
			}
		}
		// create output shape
		outputShape := Shape{
			Type:  outputType.Instantiate(bindings),
			Stack: gotInputShape.Stack,
		}
		// set new type variables on action
		funAction3 := func(inputState states.State, args []states.Action) *states.Thunk {
			for n, t := range bindings {
				inputState.TypeStack = inputState.TypeStack.Update(n, t)
			}
			return funAction2(inputState, args)
		}
		// return
		return outputShape, funAction3, ids, true, nil
	}
}
