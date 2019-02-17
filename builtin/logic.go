package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initLogic() {
	InitialShape.FuncerStack = InitialShape.FuncerStack.PushAll([]functions.Funcer{
		functions.SimpleFuncer(
			types.AnyType,
			"true",
			nil,
			types.BoolType,
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				return values.BoolValue(true)
			},
		),
		functions.SimpleFuncer(
			types.AnyType,
			"false",
			nil,
			types.BoolType,
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				return values.BoolValue(false)
			},
		),
		functions.SimpleFuncer(
			types.BoolType,
			"and",
			[]types.Type{types.BoolType},
			types.BoolType,
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				inputBool := inputValue.(values.BoolValue)
				argumentBool := argumentValues[0].(values.BoolValue)
				return values.BoolValue(inputBool && argumentBool)
			},
		),
		functions.SimpleFuncer(
			types.BoolType,
			"or",
			[]types.Type{types.BoolType},
			types.BoolType,
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				inputBool := inputValue.(values.BoolValue)
				argumentBool := argumentValues[0].(values.BoolValue)
				return values.BoolValue(inputBool || argumentBool)
			},
		),
		functions.SimpleFuncer(
			types.BoolType,
			"not",
			nil,
			types.BoolType,
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				inputBool := inputValue.(values.BoolValue)
				return values.BoolValue(!inputBool)
			},
		),
		functions.SimpleFuncer(
			types.BoolType,
			"==",
			[]types.Type{types.BoolType},
			types.BoolType,
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				inputBool := inputValue.(values.BoolValue)
				argumentBool := argumentValues[0].(values.BoolValue)
				return values.BoolValue(inputBool == argumentBool)
			},
		),
	})
}
