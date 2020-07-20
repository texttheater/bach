package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initObj() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.RegularFuncer(
			types.ObjType{
				PropTypeMap: map[string]types.Type{},
				RestType: types.TypeVariable{
					Name:       "A",
					UpperBound: types.AnyType{},
				},
			},
			"get",
			[]*states.Parameter{
				states.SimpleParam(types.Union(
					types.StrType{},
					types.NumType{},
				)),
			},
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
			func(inputState states.State, args []states.Action) *states.Thunk {
				inputValue := inputState.Value.(states.ObjValue)
				res0 := args[0](inputState, nil).Eval()
				if res0.Error != nil {
					return states.ThunkFromError(res0.Error)
				}
				prop, err := res0.Value.Out() // TODO ???
				if err != nil {
					return states.ThunkFromError(err)
				}
				thunk, ok := inputValue[prop]
				if !ok {
					return states.ThunkFromValue(states.NullValue{})
				}
				res := thunk.Eval()
				if res.Error != nil {
					return states.ThunkFromError(res.Error)
				}
				return states.ThunkFromValue(states.ObjValue(
					map[string]*states.Thunk{
						"just": states.ThunkFromValue(res.Value),
					},
				))
			},
			nil,
		),
	})
}
