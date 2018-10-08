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
