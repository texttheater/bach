package builtin

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initControl() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.Funcer{
		expressions.RegularFuncer(
			types.AnyType{},
			"fatal",
			nil,
			types.VoidType{},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				return states.ThunkFromError(
					errors.E(
						errors.Code(errors.UnexpectedValue),
						errors.Pos(pos),
						errors.GotValue(inputState.Value),
					),
				)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.Union(
				types.NullType{},
				types.TypeVariable{
					Name:       "A",
					UpperBound: types.AnyType{},
				},
			),
			"must",
			nil,
			types.TypeVariable{
				Name:       "A",
				UpperBound: types.AnyType{},
			},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				switch inputState.Value.(type) {
				case states.NullValue:
					return states.ThunkFromError(
						errors.E(
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
	})
}
