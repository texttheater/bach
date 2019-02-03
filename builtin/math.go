package builtin

import (
	"math"

	"github.com/texttheater/bach/values"
)

func Add(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(values.NumValue)
	argumentNum := argumentValues[0].(values.NumValue)
	return values.NumValue(inputNum + argumentNum)
}

func Subtract(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(values.NumValue)
	argumentNum := argumentValues[0].(values.NumValue)
	return values.NumValue(inputNum - argumentNum)
}

func Multiply(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(values.NumValue)
	argumentNum := argumentValues[0].(values.NumValue)
	return values.NumValue(inputNum * argumentNum)
}

func Divide(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(values.NumValue)
	argumentNum := argumentValues[0].(values.NumValue)
	return values.NumValue(inputNum / argumentNum)
}

func Modulo(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(values.NumValue)
	argumentNum := argumentValues[0].(values.NumValue)
	return values.NumValue(math.Mod(float64(inputNum), float64(argumentNum)))
}

func LessThan(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(values.NumValue)
	argumentNum := argumentValues[0].(values.NumValue)
	return values.BoolValue(inputNum < argumentNum)
}

func GreaterThan(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(values.NumValue)
	argumentNum := argumentValues[0].(values.NumValue)
	return values.BoolValue(inputNum > argumentNum)
}

func Equal(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(values.NumValue)
	argumentNum := argumentValues[0].(values.NumValue)
	return values.BoolValue(inputNum == argumentNum)
}

func LessEqual(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(values.NumValue)
	argumentNum := argumentValues[0].(values.NumValue)
	return values.BoolValue(inputNum <= argumentNum)
}

func GreaterEqual(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputNum := inputValue.(values.NumValue)
	argumentNum := argumentValues[0].(values.NumValue)
	return values.BoolValue(inputNum >= argumentNum)
}
