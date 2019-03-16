package builtin

import (
	"math"

	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initMath() {
	InitialShape.FuncerStack = InitialShape.FuncerStack.PushAll([]shapes.Funcer{
		shapes.SimpleFuncer(
			types.NumType,
			"+",
			[]types.Type{types.NumType},
			types.NumType,
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.NumValue(inputNum + argumentNum)
			},
		),
		shapes.SimpleFuncer(
			types.NumType,
			"-",
			[]types.Type{types.NumType},
			types.NumType,
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.NumValue(inputNum - argumentNum)
			},
		),
		shapes.SimpleFuncer(
			types.NumType,
			"*",
			[]types.Type{types.NumType},
			types.NumType,
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.NumValue(inputNum * argumentNum)
			},
		),
		shapes.SimpleFuncer(
			types.NumType,
			"/",
			[]types.Type{types.NumType},
			types.NumType,
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.NumValue(inputNum / argumentNum)
			},
		),
		shapes.SimpleFuncer(
			types.NumType,
			"%",
			[]types.Type{types.NumType},
			types.NumType,
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.NumValue(math.Mod(float64(inputNum), float64(argumentNum)))
			},
		),
		shapes.SimpleFuncer(
			types.NumType,
			"<",
			[]types.Type{types.NumType},
			types.BoolType,
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.BoolValue(inputNum < argumentNum)
			},
		),
		shapes.SimpleFuncer(
			types.NumType,
			">",
			[]types.Type{types.NumType},
			types.BoolType,
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.BoolValue(inputNum > argumentNum)
			},
		),
		shapes.SimpleFuncer(
			types.NumType,
			"==",
			[]types.Type{types.NumType},
			types.BoolType,
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.BoolValue(inputNum == argumentNum)
			},
		),
		shapes.SimpleFuncer(
			types.NumType,
			"<=",
			[]types.Type{types.NumType},
			types.BoolType,
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.BoolValue(inputNum <= argumentNum)
			},
		),
		shapes.SimpleFuncer(
			types.NumType,
			">=",
			[]types.Type{types.NumType},
			types.BoolType,
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				inputNum := inputValue.(values.NumValue)
				argumentNum := argumentValues[0].(values.NumValue)
				return values.BoolValue(inputNum >= argumentNum)
			},
		),
	})
}
