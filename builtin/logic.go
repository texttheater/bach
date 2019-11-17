package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initLogic() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.SimpleFuncer(
			types.AnyType{},
			"true",
			nil,
			types.BoolType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				return values.BoolValue(true), nil
			},
		),
		functions.SimpleFuncer(
			types.AnyType{},
			"false",
			nil,
			types.BoolType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				return values.BoolValue(false), nil
			},
		),
		functions.SimpleFuncer(
			types.BoolType{},
			"and",
			[]types.Type{types.BoolType{}},
			types.BoolType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				inputBool := inputValue.(values.BoolValue)
				argumentBool := argumentValues[0].(values.BoolValue)
				return values.BoolValue(inputBool && argumentBool), nil
			},
		),
		functions.SimpleFuncer(
			types.BoolType{},
			"or",
			[]types.Type{types.BoolType{}},
			types.BoolType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				inputBool := inputValue.(values.BoolValue)
				argumentBool := argumentValues[0].(values.BoolValue)
				return values.BoolValue(inputBool || argumentBool), nil
			},
		),
		functions.SimpleFuncer(
			types.BoolType{},
			"not",
			nil,
			types.BoolType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				inputBool := inputValue.(values.BoolValue)
				return values.BoolValue(!inputBool), nil
			},
		),
		functions.SimpleFuncer(
			types.BoolType{},
			"==",
			[]types.Type{types.BoolType{}},
			types.BoolType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				inputBool := inputValue.(values.BoolValue)
				argumentBool := argumentValues[0].(values.BoolValue)
				return values.BoolValue(inputBool == argumentBool), nil
			},
		),
	})
}
