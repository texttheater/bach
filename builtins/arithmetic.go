package builtins

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

///////////////////////////////////////////////////////////////////////////////

type Add struct {
	Arg functions.Function
}

func (f Add) Type() types.Type {
	return &types.NumberType{}
}

func (f Add) Evaluate(input values.Value) values.Value {
	numberInput, _ := input.(*values.NumberValue)
	argValue := f.Arg.Evaluate(&values.NullValue{})
	numberArgValue, _ := argValue.(*values.NumberValue)
	return &values.NumberValue{numberInput.Value + numberArgValue.Value}
}

///////////////////////////////////////////////////////////////////////////////

type Subtract struct {
	Arg functions.Function
}

func (f Subtract) Type() types.Type {
	return &types.NumberType{}
}

func (f Subtract) Evaluate(input values.Value) values.Value {
	numberInput, _ := input.(*values.NumberValue)
	argValue := f.Arg.Evaluate(&values.NullValue{})
	numberArgValue, _ := argValue.(*values.NumberValue)
	return &values.NumberValue{numberInput.Value - numberArgValue.Value}
}

///////////////////////////////////////////////////////////////////////////////

type Multiply struct {
	Arg functions.Function
}

func (f Multiply) Type() types.Type {
	return &types.NumberType{}
}

func (f Multiply) Evaluate(input values.Value) values.Value {
	numberInput, _ := input.(*values.NumberValue)
	argValue := f.Arg.Evaluate(&values.NullValue{})
	numberArgValue, _ := argValue.(*values.NumberValue)
	return &values.NumberValue{numberInput.Value * numberArgValue.Value}
}

///////////////////////////////////////////////////////////////////////////////

type Divide struct {
	Arg functions.Function
}

func (f Divide) Type() types.Type {
	return &types.NumberType{}
}

func (f Divide) Evaluate(input values.Value) values.Value {
	numberInput, _ := input.(*values.NumberValue)
	argValue := f.Arg.Evaluate(&values.NullValue{})
	numberArgValue, _ := argValue.(*values.NumberValue)
	return &values.NumberValue{numberInput.Value / numberArgValue.Value}
}
