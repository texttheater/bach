package builtin

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var ControlFuncers = []shapes.Funcer{
	shapes.Funcer{
		Summary:           "Aborts with an error message",
		InputType:         types.Any{},
		InputDescription:  "any value",
		Name:              "fatal",
		Params:            nil,
		OutputType:        types.Void{},
		OutputDescription: "does not return",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return states.ThunkFromError(
				errors.ValueError(
					errors.Code(errors.UnexpectedValue),
					errors.Pos(pos),
					errors.GotValue(inputState.Value),
				),
			)
		},
		Examples: []shapes.Example{
			{`1 if ==1 then fatal else true ok`, `Bool`, ``, errors.ValueError(
				errors.Code(errors.UnexpectedValue),
				errors.Pos(lexer.Position{Offset: 14, Line: 1, Column: 15}),
				errors.GotValue(states.NumValue(1)),
			)},
		},
		IDs: nil,
	},
	shapes.Funcer{
		Summary: "Aborts with an error message if input is null",
		InputType: types.NewUnion(
			types.Null{},
			types.NewVar("A", types.Any{}),
		),
		InputDescription:  "a value that you want to be sure is not null",
		Name:              "must",
		Params:            nil,
		OutputType:        types.NewVar("A", types.Any{}),
		OutputDescription: "the input value, unless it is null",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			switch inputState.Value.(type) {
			case states.NullValue:
				return states.ThunkFromError(
					errors.ValueError(
						errors.Code(errors.UnexpectedValue),
						errors.Pos(pos),
						errors.GotValue(inputState.Value),
					),
				)
			default:
				return states.ThunkFromValue(inputState.Value)
			}
		},
		Examples: []shapes.Example{
			{`null must`, `Null`, ``, errors.ValueError(
				errors.Code(errors.UnexpectedValue),
				errors.Pos(lexer.Position{Offset: 5, Line: 1, Column: 6}),
				errors.GotValue(states.NullValue{}),
			)},
		},
		IDs: nil,
	},
}
