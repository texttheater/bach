package builtin

import (
	"math"

	"github.com/texttheater/bach/values"
)

func Add(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(*values.NumValue)
	argumentNum := argumentValues[0].(*values.NumValue)
	return &values.NumValue{inputNum.Value + argumentNum.Value}
}

func Subtract(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(*values.NumValue)
	argumentNum := argumentValues[0].(*values.NumValue)
	return &values.NumValue{inputNum.Value - argumentNum.Value}
}

func Multiply(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(*values.NumValue)
	argumentNum := argumentValues[0].(*values.NumValue)
	return &values.NumValue{inputNum.Value * argumentNum.Value}
}

func Divide(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(*values.NumValue)
	argumentNum := argumentValues[0].(*values.NumValue)
	return &values.NumValue{inputNum.Value / argumentNum.Value}
}

func Modulo(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(*values.NumValue)
	argumentNum := argumentValues[0].(*values.NumValue)
	return &values.NumValue{math.Mod(inputNum.Value, argumentNum.Value)}
}

func LessThan(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(*values.NumValue)
	argumentNum := argumentValues[0].(*values.NumValue)
	return &values.BoolValue{inputNum.Value < argumentNum.Value}
}

func GreaterThan(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(*values.NumValue)
	argumentNum := argumentValues[0].(*values.NumValue)
	return &values.BoolValue{inputNum.Value > argumentNum.Value}
}

func Equal(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(*values.NumValue)
	argumentNum := argumentValues[0].(*values.NumValue)
	return &values.BoolValue{inputNum.Value == argumentNum.Value}
}

func LessEqual(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(*values.NumValue)
	argumentNum := argumentValues[0].(*values.NumValue)
	return &values.BoolValue{inputNum.Value <= argumentNum.Value}
}

func GreaterEqual(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(*values.NumValue)
	argumentNum := argumentValues[0].(*values.NumValue)
	return &values.BoolValue{inputNum.Value >= argumentNum.Value}
}
