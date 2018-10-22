package builtin

import (
	"math"

	"github.com/texttheater/bach/values"
)

func Add(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNumber := inputValue.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.NumberValue{inputNumber.Value + argumentNumber.Value}
}

func Subtract(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNumber := inputValue.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.NumberValue{inputNumber.Value - argumentNumber.Value}
}

func Multiply(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNumber := inputValue.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.NumberValue{inputNumber.Value * argumentNumber.Value}
}

func Divide(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNumber := inputValue.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.NumberValue{inputNumber.Value / argumentNumber.Value}
}

func Modulo(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNumber := inputValue.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.NumberValue{math.Mod(inputNumber.Value, argumentNumber.Value)}
}

func LessThan(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNumber := inputValue.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.BooleanValue{inputNumber.Value < argumentNumber.Value}
}

func GreaterThan(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNumber := inputValue.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.BooleanValue{inputNumber.Value > argumentNumber.Value}
}

func Equal(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNumber := inputValue.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.BooleanValue{inputNumber.Value == argumentNumber.Value}
}

func LessEqual(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNumber := inputValue.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.BooleanValue{inputNumber.Value <= argumentNumber.Value}
}

func GreaterEqual(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNumber := inputValue.(*values.NumberValue)
	argumentNumber := argumentValues[0].(*values.NumberValue)
	return &values.BooleanValue{inputNumber.Value >= argumentNumber.Value}
}
