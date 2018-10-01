package functions

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

///////////////////////////////////////////////////////////////////////////////

type Function interface {
	Type() types.Type
	Evaluate(input values.Value) values.Value
}

///////////////////////////////////////////////////////////////////////////////

type IdentityFunction struct {
	Type_ types.Type
}

func (f IdentityFunction) Type() types.Type {
	return f.Type_
}

func (f IdentityFunction) Evaluate(input values.Value) values.Value {
	return input
}

///////////////////////////////////////////////////////////////////////////////

type CompositionFunction struct {
	Left  Function
	Right Function
}

func (f CompositionFunction) Type() types.Type {
	return f.Right.Type()
}

func (f CompositionFunction) Evaluate(input values.Value) values.Value {
	return f.Right.Evaluate(f.Left.Evaluate(input))
}

///////////////////////////////////////////////////////////////////////////////

type NumberFunction struct {
	Value float64
}

func (f NumberFunction) Type() types.Type {
	return types.NumberType{}
}

func (f NumberFunction) Evaluate(input values.Value) values.Value {
	return values.NumberValue{f.Value}
}
