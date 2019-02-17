package builtin

import (
	"math"

	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initMath() {
	InitialShape.FuncerStack = InitialShape.FuncerStack.PushAll([]functions.Funcer{
		functions.SimpleFuncer(
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
		functions.SimpleFuncer(
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
		functions.SimpleFuncer(
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
		functions.SimpleFuncer(
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
		functions.SimpleFuncer(
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
		functions.SimpleFuncer(
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
		functions.SimpleFuncer(
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
		functions.SimpleFuncer(
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
		functions.SimpleFuncer(
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
		functions.SimpleFuncer(
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
