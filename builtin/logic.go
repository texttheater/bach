package builtin

import (
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initLogic() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]shapes.Funcer{
		shapes.SimpleFuncer(
			types.AnyType{},
			"true",
			nil,
			types.BoolType{},
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				return values.BoolValue(true)
			},
		),
		shapes.SimpleFuncer(
			types.AnyType{},
			"false",
			nil,
			types.BoolType{},
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				return values.BoolValue(false)
			},
		),
		shapes.SimpleFuncer(
			types.BoolType{},
			"and",
			[]types.Type{types.BoolType{}},
			types.BoolType{},
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				inputBool := inputValue.(values.BoolValue)
				argumentBool := argumentValues[0].(values.BoolValue)
				return values.BoolValue(inputBool && argumentBool)
			},
		),
		shapes.SimpleFuncer(
			types.BoolType{},
			"or",
			[]types.Type{types.BoolType{}},
			types.BoolType{},
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				inputBool := inputValue.(values.BoolValue)
				argumentBool := argumentValues[0].(values.BoolValue)
				return values.BoolValue(inputBool || argumentBool)
			},
		),
		shapes.SimpleFuncer(
			types.BoolType{},
			"not",
			nil,
			types.BoolType{},
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				inputBool := inputValue.(values.BoolValue)
				return values.BoolValue(!inputBool)
			},
		),
		shapes.SimpleFuncer(
			types.BoolType{},
			"==",
			[]types.Type{types.BoolType{}},
			types.BoolType{},
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				inputBool := inputValue.(values.BoolValue)
				argumentBool := argumentValues[0].(values.BoolValue)
				return values.BoolValue(inputBool == argumentBool)
			},
		),
	})
}
