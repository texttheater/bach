package shapes

import (
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

///////////////////////////////////////////////////////////////////////////////

type IdentityFunction struct {
}

func (f *IdentityFunction) OutputShape(inputShape Shape) Shape {
	return inputShape
}

func (f *IdentityFunction) OutputState(inputState states.State) states.State {
	return inputState
}

///////////////////////////////////////////////////////////////////////////////

type CompositionFunction struct {
	Left Function
	Right Function
}

func (f *CompositionFunction) OutputShape(inputShape Shape) Shape {
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

func (f *NumberFunction) OutputShape(inputShape Shape) Shape {
	return Shape{&types.NumberType{}, inputShape.Stack}
}

func (f *NumberFunction) OutputState(inputState states.State) states.State {
	return states.State{&values.NumberValue{f.Value}, inputState.Stack}
}

///////////////////////////////////////////////////////////////////////////////
