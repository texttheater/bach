package functions

import (
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

///////////////////////////////////////////////////////////////////////////////

type IdentityFunction struct {
}

func (f *IdentityFunction) OutputShape(inputShape shapes.Shape) shapes.Shape {
	return inputShape
}

func (f *IdentityFunction) OutputState(inputState states.State) states.State {
	return inputState
}

///////////////////////////////////////////////////////////////////////////////

type CompositionFunction struct {
	Left  shapes.Function
	Right shapes.Function
}

func (f *CompositionFunction) OutputShape(inputShape shapes.Shape) shapes.Shape {
	return f.Right.OutputShape(f.Left.OutputShape(inputShape))
}

func (f *CompositionFunction) OutputState(inputState states.State) states.State {
	return f.Right.OutputState(f.Left.OutputState(inputState))
}

///////////////////////////////////////////////////////////////////////////////

type LiteralFunction struct {
	Type  types.Type
	Value values.Value
}

func (f *LiteralFunction) OutputShape(inputShape shapes.Shape) shapes.Shape {
	return shapes.Shape{f.Type, inputShape.Stack}
}

func (f *LiteralFunction) OutputState(inputState states.State) states.State {
	return states.State{f.Value, inputState.Stack}
}

///////////////////////////////////////////////////////////////////////////////

type EvaluatorFunction struct {
	ArgumentFunctions []shapes.Function
	OutputType        types.Type
	Kernel            func(inputValue values.Value, argumentValues []values.Value) values.Value
}

func (f *EvaluatorFunction) OutputShape(inputShape shapes.Shape) shapes.Shape {
	return shapes.Shape{f.OutputType, inputShape.Stack}
}

func (f *EvaluatorFunction) OutputState(inputState states.State) states.State {
	argumentInput := states.State{&values.NullValue{}, inputState.Stack}
	argumentValues := make([]values.Value, len(f.ArgumentFunctions))
	for i, a := range f.ArgumentFunctions {
		argumentValues[i] = a.OutputState(argumentInput).Value
	}
	return states.State{f.Kernel(inputState.Value, argumentValues), inputState.Stack}
}

///////////////////////////////////////////////////////////////////////////////

type AssignmentFunction struct {
	Name string
}

func (f *AssignmentFunction) OutputShape(inputShape shapes.Shape) shapes.Shape {
	return shapes.Shape{inputShape.Type, inputShape.Stack.Push(shapes.NFF{
		&types.AnyType{},
		f.Name,
		[]types.Type{},
		func(argumentFunctions []shapes.Function) shapes.Function {
			return &VariableFunction{f.Name, inputShape.Type}
		},
	})}
}

func (f *AssignmentFunction) OutputState(inputState states.State) states.State {
	return states.State{
		inputState.Value,
		inputState.Stack.Push(states.NamedValue{
			f.Name,
			inputState.Value,
		}),
	}
}

///////////////////////////////////////////////////////////////////////////////

type VariableFunction struct {
	Name string
	Type types.Type
}

func (f *VariableFunction) OutputShape(inputShape shapes.Shape) shapes.Shape {
	return shapes.Shape{f.Type, inputShape.Stack}
}

func (f *VariableFunction) OutputState(inputState states.State) states.State {
	stack := inputState.Stack
	for stack != nil {
		if stack.Head.Name == f.Name {
			return states.State{
				stack.Head.Value,
				inputState.Stack,
			}
		}
		stack = stack.Tail
	}
	panic("variable not found")
}

///////////////////////////////////////////////////////////////////////////////
