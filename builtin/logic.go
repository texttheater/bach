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
		types.Any{},
		"any value (is ignored)",
		"true",
		nil,
		types.Bool{},
		"the unique value representing logical truth",
		"",
		func(inputThunk *states.Thunk, argumentThunks []*states.Thunk) *states.Thunk {
			return states.ThunkFromValue(states.BoolValue(true))
		},
		[]shapes.Example{
			{`true`, `Bool`, `true`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Returns the value representing logical falsehood.",
		types.Any{},
		"any value (is ignored)",
		"false",
		nil,
		types.Bool{},
		"the unique value representing logical falsehood",
		"",
		func(inputThunk *states.Thunk, argumentThunks []*states.Thunk) *states.Thunk {
			return states.ThunkFromValue(states.BoolValue(false))
		},
		[]shapes.Example{
			{`false`, `Bool`, `false`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Computes logical conjunction.",
		types.Bool{},
		"a boolean value",
		"and",
		[]*params.Param{
			params.SimpleParam("q", "another boolean value", types.Bool{}),
		},
		types.Bool{},
		"true if both the input and q are true, false otherwise",
		"",
		func(inputThunk *states.Thunk, argumentThunks []*states.Thunk) *states.Thunk {
			inputBool, err := inputThunk.EvalBool()
			if err != nil {
				return states.ThunkFromError(err)
			}
			argumentBool, err := argumentThunks[0].EvalBool()
			if err != nil {
				return states.ThunkFromError(err)
			}
			return states.ThunkFromValue(states.BoolValue(inputBool && argumentBool))
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
		types.Bool{},
		"a boolean value",
		"or",
		[]*params.Param{
			params.SimpleParam("q", "another boolean value", types.Bool{}),
		},
		types.Bool{},
		"true if at least one of the input and q is true, false otherwise",
		"",
		func(inputThunk *states.Thunk, argumentThunks []*states.Thunk) *states.Thunk {
			inputBool, err := inputThunk.EvalBool()
			if err != nil {
				return states.ThunkFromError(err)
			}
			argumentBool, err := argumentThunks[0].EvalBool()
			if err != nil {
				return states.ThunkFromError(err)
			}
			return states.ThunkFromValue(states.BoolValue(inputBool || argumentBool))
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
		types.Bool{},
		"a boolean value",
		"not",
		nil,
		types.Bool{},
		"true if the input is false, and false if the input is true",
		"",
		func(inputThunk *states.Thunk, argumentThunks []*states.Thunk) *states.Thunk {
			inputBool, err := inputThunk.EvalBool()
			if err != nil {
				return states.ThunkFromError(err)
			}
			return states.ThunkFromValue(states.BoolValue(!inputBool))
		},
		[]shapes.Example{
			{`false not`, `Bool`, `true`, nil},
			{`true not`, `Bool`, `false`, nil},
			{`1 +1 ==2 and(2 +2 ==5 not)`, `Bool`, `true`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Checks equality of boolean values.",
		types.Bool{},
		"a boolean value",
		"==",
		[]*params.Param{
			params.SimpleParam("q", "another boolean value", types.Bool{}),
		},
		types.Bool{},
		"true if the input and q are identical, false otherwise",
		"",
		func(inputThunk *states.Thunk, argumentThunks []*states.Thunk) *states.Thunk {
			inputBool, err := inputThunk.EvalBool()
			if err != nil {
				return states.ThunkFromError(err)
			}
			argumentBool, err := argumentThunks[0].EvalBool()
			if err != nil {
				return states.ThunkFromError(err)
			}
			return states.ThunkFromValue(states.BoolValue(inputBool == argumentBool))
		},
		[]shapes.Example{
			{`false ==false`, `Bool`, `true`, nil},
			{`false ==true`, `Bool`, `false`, nil},
			{`true ==false`, `Bool`, `false`, nil},
			{`true ==true`, `Bool`, `true`, nil},
		},
	),
}
