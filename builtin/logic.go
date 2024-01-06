package builtin

import (
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var LogicFuncers = []expressions.Funcer{
	expressions.SimpleFuncer(
		types.Any{},
		"true",
		nil,
		types.Bool{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.BoolValue(true), nil
		},
	),
	expressions.SimpleFuncer(
		types.Any{},
		"false",
		nil,
		types.Bool{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.BoolValue(false), nil
		},
	),
	expressions.SimpleFuncer(
		types.Bool{},
		"and",
		[]*params.Param{
			params.SimpleParam("q", types.Bool{}),
		},
		types.Bool{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			inputBool := inputValue.(states.BoolValue)
			argumentBool := argumentValues[0].(states.BoolValue)
			return states.BoolValue(inputBool && argumentBool), nil
		},
	),
	expressions.SimpleFuncer(
		types.Bool{},
		"or",
		[]*params.Param{
			params.SimpleParam("q", types.Bool{}),
		},
		types.Bool{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			inputBool := inputValue.(states.BoolValue)
			argumentBool := argumentValues[0].(states.BoolValue)
			return states.BoolValue(inputBool || argumentBool), nil
		},
	),
	expressions.SimpleFuncer(
		types.Bool{},
		"not",
		nil,
		types.Bool{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			inputBool := inputValue.(states.BoolValue)
			return states.BoolValue(!inputBool), nil
		},
	),
	expressions.SimpleFuncer(
		types.Bool{},
		"==",
		[]*params.Param{
			params.SimpleParam("q", types.Bool{}),
		},
		types.Bool{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			inputBool := inputValue.(states.BoolValue)
			argumentBool := argumentValues[0].(states.BoolValue)
			return states.BoolValue(inputBool == argumentBool), nil
		},
	),
}
