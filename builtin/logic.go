package builtin

import (
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var LogicFuncers = []shapes.Funcer{
	shapes.SimpleFuncer(
		"Returns the value representing logical truth.",
		types.AnyType{},
		"any value (is ignored)",
		"true",
		nil,
		types.BoolType{},
		"the unique value representing logical truth",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.BoolValue(true), nil
		},
		[]shapes.Example{
			{`true`, `Bool`, `true`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Returns the value representing logical falsehood.",
		types.AnyType{},
		"any value (is ignored)",
		"false",
		nil,
		types.BoolType{},
		"the unique value representing logical falsehood",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.BoolValue(false), nil
		},
		[]shapes.Example{
			{`false`, `Bool`, `false`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes logical conjunction.",
		types.BoolType{},
		"a boolean value",
		"and",
		[]*params.Param{
			params.SimpleParam("q", "another boolean value", types.BoolType{}),
		},
		types.BoolType{},
		"true if both the input and q are true, false otherwise",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			inputBool := inputValue.(states.BoolValue)
			argumentBool := argumentValues[0].(states.BoolValue)
			return states.BoolValue(inputBool && argumentBool), nil
		},
		[]shapes.Example{
			{`false and(false)`, `Bool`, `false`, nil},
			{`false and(true)`, `Bool`, `false`, nil},
			{`true and(false)`, `Bool`, `false`, nil},
			{`true and(true)`, `Bool`, `true`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes logical disjunction.",
		types.BoolType{},
		"a boolean value",
		"or",
		[]*params.Param{
			params.SimpleParam("q", "another boolean value", types.BoolType{}),
		},
		types.BoolType{},
		"true if at least one of the input and q is true, false otherwise",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			inputBool := inputValue.(states.BoolValue)
			argumentBool := argumentValues[0].(states.BoolValue)
			return states.BoolValue(inputBool || argumentBool), nil
		},
		[]shapes.Example{
			{`false or(false)`, `Bool`, `false`, nil},
			{`false or(true)`, `Bool`, `true`, nil},
			{`true or(false)`, `Bool`, `true`, nil},
			{`true or(true)`, `Bool`, `true`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes logical negation.",
		types.BoolType{},
		"a boolean value",
		"not",
		nil,
		types.BoolType{},
		"true if the input is false, and false if the input is true",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			inputBool := inputValue.(states.BoolValue)
			return states.BoolValue(!inputBool), nil
		},
		[]shapes.Example{
			{`false not`, `Bool`, `true`, nil},
			{`true not`, `Bool`, `false`, nil},
			{`1 +1 ==2 and(2 +2 ==5 not)`, `Bool`, `true`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Checks equality of boolean values.",
		types.BoolType{},
		"a boolean value",
		"==",
		[]*params.Param{
			params.SimpleParam("q", "another boolean value", types.BoolType{}),
		},
		types.BoolType{},
		"true if the input and q are identical, false otherwise",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			inputBool := inputValue.(states.BoolValue)
			argumentBool := argumentValues[0].(states.BoolValue)
			return states.BoolValue(inputBool == argumentBool), nil
		},
		[]shapes.Example{
			{`false ==false`, `Bool`, `true`, nil},
			{`false ==true`, `Bool`, `false`, nil},
			{`true ==false`, `Bool`, `false`, nil},
			{`true ==true`, `Bool`, `true`, nil},
		},
	),
}
