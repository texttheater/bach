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

func (x CallExpression) Typecheck(inputShape Shape, p []*params.Param) (Shape, states.Action, *states.IDStack, error) {
	// go down the function stack and find the function invoked by this
	// call
	stack := inputShape.Stack
funcers:
	for {
		// reached bottom of stack without finding a matching funcer
		if stack == nil {
			return Shape{}, nil, nil, errors.TypeError(
				errors.Code(errors.NoSuchFunction),
				errors.Pos(x.Pos),
				errors.InputType(inputShape.Type),
				errors.Name(x.Name),
				errors.NumParams(len(x.Args)+len(p)),
			)
		}
		// try the funcer on top of the stack
		funcerDefinition := stack.Head
		ids := funcerDefinition.IDs
		// match number of parameters
		if len(x.Args)+len(p) != len(funcerDefinition.Params) {
			stack = stack.Tail
			continue funcers
		}
		// match name
		if x.Name != funcerDefinition.Name {
			stack = stack.Tail
			continue funcers
		}
		// match input type
		bindings := make(map[string]types.Type)
		if !funcerDefinition.InputType.Bind(inputShape.Type, bindings) {
			stack = stack.Tail
			continue funcers
		}
		if !funcerDefinition.InputType.Instantiate(bindings).Subsumes(inputShape.Type) {
			stack = stack.Tail
			continue funcers
		}
		// typecheck and set parameters filled by this call
		action1 := func(inputState states.State, args []states.Action) *states.Thunk {
			return funcerDefinition.Kernel(inputState, args, bindings, x.Position())
		}
		argActions := make([]states.Action, len(x.Args))
		argIDss := make([]*states.IDStack, len(x.Args))
		for i := range x.Args {
			argInputShape := Shape{
				Type: funcerDefinition.Params[i].InputType.Instantiate(bindings),
				// TODO what if we don't have the bindings yet?
				// TODO what if an incompatible bound is declared?
				Stack: inputShape.Stack,
			}
			argOutputShape, argAction, argIDs, err := x.Args[i].Typecheck(
				argInputShape, instantiate(funcerDefinition.Params[i].Params, bindings))
			if err != nil {
				stack = stack.Tail
				continue funcers
			}
			if !funcerDefinition.Params[i].OutputType.Bind(argOutputShape.Type, bindings) {
				return Shape{}, nil, nil, errors.TypeError(
					errors.Code(errors.ArgHasWrongOutputType),
					errors.Pos(x.Pos),
					errors.ArgNum(i+1),
					errors.WantType(funcerDefinition.Params[i].OutputType.Instantiate(bindings)),
					errors.GotType(argOutputShape.Type),
				)
			}
			if !funcerDefinition.Params[i].OutputType.Instantiate(bindings).Subsumes(argOutputShape.Type) {
				return Shape{}, nil, nil, errors.TypeError(
					errors.Code(errors.ArgHasWrongOutputType),
					errors.Pos(x.Pos),
					errors.ArgNum(i+1),
					errors.WantType(funcerDefinition.Params[i].OutputType.Instantiate(bindings)),
					errors.GotType(argOutputShape.Type),
				)
			}
			argActions[i] = argAction
			argIDss[i] = argIDs
			ids = ids.AddAll(argIDs)
		}
		// pass input variable stack to arguments
		action2 := func(inputState states.State, args []states.Action) *states.Thunk {
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
			return action1(inputState, args2)
		}
		// typecheck parameters not filled by the call
		for i, gotParam := range p {
			wantParam := funcerDefinition.Params[len(x.Args)+i].Instantiate(bindings)
			if !gotParam.Subsumes(wantParam) {
				return Shape{}, nil, nil, errors.TypeError(
					errors.Code(errors.ParamDoesNotMatch),
					errors.Pos(x.Pos),
					errors.ParamNum(i+1),
					errors.WantParam(gotParam),
					errors.GotParam(wantParam),
				)
			}
		}
		// create output shape
		outputShape := Shape{
			Type:  funcerDefinition.OutputType.Instantiate(bindings),
			Stack: inputShape.Stack,
		}
		// set new type variables on action
		action3 := func(inputState states.State, args []states.Action) *states.Thunk {
			for n, t := range bindings {
				inputState.TypeStack = inputState.TypeStack.Update(n, t)
			}
			return action2(inputState, args)
		}
		// make call lazy
		action := func(inputState states.State, args []states.Action) *states.Thunk {
			return &states.Thunk{
				Func: func() *states.Thunk {
					return action3(inputState, args)
				},
				Stack:     inputState.Stack,
				TypeStack: inputState.TypeStack,
			}
		}
		// return
		return outputShape, action, ids, nil
	}
}

type Shape struct {
	Type  types.Type
	Stack *FuncerStack
}

type FuncerStack struct {
	Head FuncerDefinition
	Tail *FuncerStack
}

func (s *FuncerStack) Push(funcer FuncerDefinition) *FuncerStack {
	return &FuncerStack{funcer, s}
}

func (s *FuncerStack) PushAll(funcers []FuncerDefinition) *FuncerStack {
	for i := len(funcers) - 1; i >= 0; i-- {
		s = s.Push(funcers[i])
	}
	return s
}

func (s *FuncerStack) String() string {
	slice := make([]FuncerDefinition, 0)
	stack := s
	for stack != nil {
		slice = append(slice, stack.Head)
		stack = stack.Tail
	}
	return fmt.Sprintf("%v", slice)
}

type FuncerDefinition struct {
	InputType  types.Type
	Name       string
	Params     []*params.Param
	OutputType types.Type
	Kernel     RegularKernel
	IDs        *states.IDStack
}

func instantiate(pars []*params.Param, bindings map[string]types.Type) []*params.Param {
	result := make([]*(params.Param), len(pars))
	for i, par := range pars {
		result[i] = par.Instantiate(bindings)
	}
	return result
}

type SimpleKernel func(inputValue states.Value, argValues []states.Value) (states.Value, error)

func SimpleFuncer(wantInputType types.Type, wantName string, argTypes []types.Type, outputType types.Type, simpleKernel SimpleKernel) FuncerDefinition {
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
			val, err := arg(inputState.Clear(), nil).Eval()
			if err != nil {
				return states.ThunkFromError(err)
			}
			argValues[i] = val
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

func VariableFuncer(id any, name string, varType types.Type) FuncerDefinition {
	kernel := func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		stack := inputState.Stack
		for stack != nil {
			if stack.Head.ID == id {
				val, err := stack.Head.Action(states.InitialState, nil).Eval()
				if err != nil {
					return states.ThunkFromError(err)
				}
				return states.ThunkFromState(states.State{
					Value:     val,
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

func RegularFuncer(wantInputType types.Type, wantName string, pars []*params.Param, outputType types.Type, kernel RegularKernel, ids *states.IDStack) FuncerDefinition {
	return FuncerDefinition{
		InputType:  wantInputType,
		Name:       wantName,
		Params:     pars,
		OutputType: outputType,
		Kernel:     kernel,
		IDs:        ids,
	}
}
