package builtin

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initControl() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.RegularFuncer(
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
		functions.RegularFuncer(
			types.Union(
				types.ObjType{
					PropTypeMap: map[string]types.Type{
						"just": types.TypeVariable{
							Name:       "A",
							UpperBound: types.AnyType{},
						},
					},
					RestType: types.AnyType{},
				},
				types.NullType{},
			),
			"must",
			nil,
			types.TypeVariable{
				Name:       "A",
				UpperBound: types.AnyType{},
			},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				switch v := inputState.Value.(type) {
				case states.ObjValue:
					res := v["just"].Eval()
					if res.Error != nil {
						return states.ThunkFromError(res.Error)
					}
					return states.ThunkFromValue(res.Value)
				default:
					return states.ThunkFromError(
						errors.E(
							errors.Code(errors.UnexpectedValue),
							errors.Pos(pos),
							errors.GotValue(inputState.Value),
						),
					)
				}
			},
			nil,
		),
	})
}
