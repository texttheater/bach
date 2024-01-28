package builtin

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var ControlFuncers = []shapes.Funcer{
	shapes.Funcer{InputType: types.Any{}, Name: "fatal", Params: nil, OutputType: types.Void{}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		return states.ThunkFromError(
			errors.ValueError(
				errors.Code(errors.UnexpectedValue),
				errors.Pos(pos),
				errors.GotValue(inputState.Value),
			),
		)
	}, IDs: nil},

	shapes.Funcer{InputType: types.NewUnion(
		types.Null{},
		types.NewVar("A", types.Any{}),
	), Name: "must", Params: nil, OutputType: types.NewVar("A", types.Any{}), Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
	}, IDs: nil},
}
