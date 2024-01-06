package builtin

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var ControlFuncers = []expressions.Funcer{
	expressions.RegularFuncer(
		types.Any{},
		"fatal",
		nil,
		types.Void{},
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return states.ThunkFromError(
				errors.ValueError(
					errors.Code(errors.UnexpectedValue),
					errors.Pos(pos),
					errors.GotValue(inputState.Value),
				),
			)
		},
		nil,
	),
	expressions.RegularFuncer(
		types.NewUnion(
			types.Null{},
			types.NewVar("A", types.Any{}),
		),
		"must",
		nil,
		types.NewVar("A", types.Any{}),
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		nil,
	),
}
