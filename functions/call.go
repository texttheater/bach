package functions

import (
	"bytes"
	"fmt"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type CallExpression struct {
	Pos  lexer.Position
	Name string
	Args []Expression
}

func (x CallExpression) Position() lexer.Position {
	return x.Pos
}

func (x CallExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	// go down the function stack and find the function invoked by this
	// call
	stack := inputShape.Stack
	for {
		// reached bottom of stack without finding a matching funcer
		if stack == nil {
			return Shape{}, nil, errors.E(
				errors.Code(errors.NoSuchFunction),
				errors.Pos(x.Pos),
				errors.InputType(inputShape.Type),
				errors.Name(x.Name),
				errors.NumParams(len(x.Args)+len(params)),
			)
		}
		// try the funcer on top of the stack
		funcer := stack.Head
		funOutputShape, funAction, ok, err := funcer(inputShape, x, params)
		if err != nil {
			return Shape{}, nil, err
		}
		if !ok {
			stack = stack.Tail
			continue
		}
		// return
		return funOutputShape, funAction, nil
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

type Parameter struct {
	InputType  types.Type
	Params     []*Parameter
	OutputType types.Type
}

func (p Parameter) Subsumes(other Parameter) bool {
	if len(p.Params) != len(other.Params) {
		return false
	}
	if !other.InputType.Subsumes(p.InputType) {
		return false
	}
	if !p.OutputType.Subsumes(other.OutputType) {
		return false
	}
	for i, otherParam := range other.Params {
		if !otherParam.Subsumes(*p.Params[i]) {
			return false
		}
	}
	return true
}

func (p Parameter) String() string {
	buffer := bytes.Buffer{}
	if !p.InputType.Subsumes(types.AnyType{}) || len(p.Params) > 0 {
		buffer.WriteString("for ")
		buffer.WriteString(p.InputType.String())
		buffer.WriteString(" ")
	}
	if len(p.Params) > 0 {
		buffer.WriteString("(")
		buffer.WriteString(p.Params[0].String())
		for _, param := range p.Params[1:] {
			buffer.WriteString(",")
			buffer.WriteString(param.String())
		}
		buffer.WriteString(")")
	}
	buffer.WriteString(" ")
	buffer.WriteString(p.OutputType.String())
	return buffer.String()
}

type Funcer func(gotInputShape Shape, gotCall CallExpression, gotParams []*Parameter) (outputShape Shape, action states.Action, ok bool, err error)

type Kernel func(inputValue values.Value, argValues []values.Value) values.Value

func SimpleFuncer(wantInputType types.Type, wantName string, argTypes []types.Type, outputType types.Type, kernel Kernel) Funcer {
	// make parameters from argument types
	params := make([]*Parameter, len(argTypes))
	for i, argType := range argTypes {
		params[i] = &Parameter{
			InputType:  types.AnyType{},
			Params:     nil,
			OutputType: argType,
		}
	}
	// make action from kernel
	action := func(inputState states.State, args []states.Action) states.State {
		argValues := make([]values.Value, len(argTypes))
		argInputState := states.State{
			Value:     &values.NullValue{},
			Stack:     inputState.Stack,
			TypeStack: inputState.TypeStack,
		}
		for i, arg := range args {
			argValues[i] = arg(argInputState, nil).Value
		}
		return states.State{
			Value:     kernel(inputState.Value, argValues),
			Stack:     inputState.Stack,
			TypeStack: inputState.TypeStack,
		}
	}
	// return
	return RegularFuncer(wantInputType, wantName, params, outputType, action)
}

func VariableFuncer(id interface{}, name string, varType types.Type) Funcer {
	varAction := func(inputState states.State, args []states.Action) states.State {
		stack := inputState.Stack
		for stack != nil {
			if stack.Head.ID == id {
				return states.State{
					Value:     stack.Head.Action(states.InitialState, nil).Value,
					Stack:     inputState.Stack,
					TypeStack: inputState.TypeStack,
				}
			}
			stack = stack.Tail
		}
		panic(fmt.Sprintf("variable %s not found", name))
	}
	return RegularFuncer(types.AnyType{}, name, nil, varType, varAction)
}

func RegularFuncer(wantInputType types.Type, wantName string, params []*Parameter, outputType types.Type, action states.Action) Funcer {
	return func(gotInputShape Shape, gotCall CallExpression, gotParams []*Parameter) (Shape, states.Action, bool, error) {
		// match number of parameters
		if len(gotCall.Args)+len(gotParams) != len(params) {
			return Shape{}, nil, false, nil
		}
		// match name
		if gotCall.Name != wantName {
			return Shape{}, nil, false, nil
		}
		// match input type
		if !wantInputType.Subsumes(gotInputShape.Type) {
			return Shape{}, nil, false, nil
		}
		// typecheck and set parameters filled by this call
		funAction := action
		for i := range gotCall.Args {
			argInputShape := Shape{
				Type:  params[i].InputType,
				Stack: gotInputShape.Stack,
			}
			argOutputShape, argAction, err := gotCall.Args[i].Typecheck(argInputShape, params[i].Params)
			if err != nil {
				return Shape{}, nil, false, err
			}
			if !params[i].OutputType.Subsumes(argOutputShape.Type) {
				return Shape{}, nil, false, errors.E(
					errors.Code(errors.ArgHasWrongOutputType),
					errors.Pos(gotCall.Pos),
					errors.ArgNum(i),
					errors.WantType(params[i].OutputType),
					errors.GotType(argOutputShape.Type),
				)
			}
			funAction = funAction.SetArg(argAction)
		}
		// typecheck parameters not filled by the call
		for i := range gotParams {
			if !gotParams[i].Subsumes(*params[len(gotCall.Args)+i]) {
				return Shape{}, nil, false, errors.E(
					errors.Code(errors.ParamDoesNotMatch),
					errors.Pos(gotCall.Pos),
					errors.ParamNum(i),
					errors.WantParam(gotParams[i]),
					errors.GotParam(params[len(gotCall.Args)+i]),
				)
			}
		}
		// create output shape
		outputShape := Shape{
			Type:  outputType,
			Stack: gotInputShape.Stack,
		}
		// return
		return outputShape, funAction, true, nil
	}
}
