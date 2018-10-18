package builtin

import (
	"math"

	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/values"
)

func Add(inputState functions.State, argumentValues []values.Value) values.Value {
	inputNumber := inputState.Value.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.NumberValue{inputNumber.Value + argumentNumber.Value}
}

func Subtract(inputState functions.State, argumentValues []values.Value) values.Value {
	inputNumber := inputState.Value.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.NumberValue{inputNumber.Value - argumentNumber.Value}
}

func Multiply(inputState functions.State, argumentValues []values.Value) values.Value {
	inputNumber := inputState.Value.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.NumberValue{inputNumber.Value * argumentNumber.Value}
}

func Divide(inputState functions.State, argumentValues []values.Value) values.Value {
	inputNumber := inputState.Value.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.NumberValue{inputNumber.Value / argumentNumber.Value}
}

func Modulo(inputState functions.State, argumentValues []values.Value) values.Value {
	inputNumber := inputState.Value.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.NumberValue{math.Mod(inputNumber.Value, argumentNumber.Value)}
}

func LessThan(inputState functions.State, argumentValues []values.Value) values.Value {
	inputNumber := inputState.Value.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.BooleanValue{inputNumber.Value < argumentNumber.Value}
}

func GreaterThan(inputState functions.State, argumentValues []values.Value) values.Value {
	inputNumber := inputState.Value.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.BooleanValue{inputNumber.Value > argumentNumber.Value}
}

func Equal(inputState functions.State, argumentValues []values.Value) values.Value {
	inputNumber := inputState.Value.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.BooleanValue{inputNumber.Value == argumentNumber.Value}
}

func LessEqual(inputState functions.State, argumentValues []values.Value) values.Value {
	inputNumber := inputState.Value.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.BooleanValue{inputNumber.Value <= argumentNumber.Value}
}

func GreaterEqual(inputState functions.State, argumentValues []values.Value) values.Value {
	inputNumber := inputState.Value.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.BooleanValue{inputNumber.Value >= argumentNumber.Value}
}
