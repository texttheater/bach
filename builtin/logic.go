package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initLogic() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.SimpleFuncer(
			types.AnyType{},
			"true",
			nil,
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.BoolValue(true), nil
			},
		),
		functions.SimpleFuncer(
			types.AnyType{},
			"false",
			nil,
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.BoolValue(false), nil
			},
		),
		functions.SimpleFuncer(
			types.BoolType{},
			"and",
			[]types.Type{types.BoolType{}},
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputBool := inputValue.(states.BoolValue)
				argumentBool := argumentValues[0].(states.BoolValue)
				return states.BoolValue(inputBool && argumentBool), nil
			},
		),
		functions.SimpleFuncer(
			types.BoolType{},
			"or",
			[]types.Type{types.BoolType{}},
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputBool := inputValue.(states.BoolValue)
				argumentBool := argumentValues[0].(states.BoolValue)
				return states.BoolValue(inputBool || argumentBool), nil
			},
		),
		functions.SimpleFuncer(
			types.BoolType{},
			"not",
			nil,
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputBool := inputValue.(states.BoolValue)
				return states.BoolValue(!inputBool), nil
			},
		),
		functions.SimpleFuncer(
			types.BoolType{},
			"==",
			[]types.Type{types.BoolType{}},
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputBool := inputValue.(states.BoolValue)
				argumentBool := argumentValues[0].(states.BoolValue)
				return states.BoolValue(inputBool == argumentBool), nil
			},
		),
	})
}
