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

// TODO generalize to LiteralFunction, storing Type and Value

type NumberFunction struct {
	Value float64
}

func (f *NumberFunction) OutputShape(inputShape shapes.Shape) shapes.Shape {
	return shapes.Shape{&types.NumberType{}, inputShape.Stack}
}

func (f *NumberFunction) OutputState(inputState states.State) states.State {
	return states.State{&values.NumberValue{f.Value}, inputState.Stack}
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
