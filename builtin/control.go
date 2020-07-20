package builtin

import (
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initControl() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		func(gotInputShape functions.Shape, gotCall functions.CallExpression, gotParams []*parameters.Parameter) (functions.Shape, states.Action, *states.IDStack, bool, error) {
			if len(gotParams) != 0 {
				return functions.Shape{}, nil, nil, false, nil
			}
			if len(gotCall.Args) != 0 {
				return functions.Shape{}, nil, nil, false, nil
			}
			if gotCall.Name != "fatal" {
				return functions.Shape{}, nil, nil, false, nil
			}
			outputShape := functions.Shape{
				Type:  types.VoidType{},
				Stack: gotInputShape.Stack,
			}
			action := func(inputState states.State, args []states.Action) *states.Thunk {
				return states.ThunkFromError(
					errors.E(
						errors.Code(errors.UnexpectedValue),
						errors.Pos(gotCall.Pos),
						errors.GotType(gotInputShape.Type),
						errors.GotValue(inputState.Value)),
				)
			}
			return outputShape, action, nil, true, nil
		},
		func(gotInputShape functions.Shape, gotCall functions.CallExpression, gotParams []*parameters.Parameter) (functions.Shape, states.Action, *states.IDStack, bool, error) {
			typeVar := types.TypeVariable{
				Name:       "A",
				UpperBound: types.AnyType{},
			}
			wantInputType := types.Union(
				types.ObjType{
					PropTypeMap: map[string]types.Type{
						"just": typeVar,
					},
					RestType: types.AnyType{},
				},
				types.NullType{},
			)
			bindings := make(map[string]types.Type)
			if !wantInputType.Bind(gotInputShape.Type, bindings) {
				return functions.Shape{}, nil, nil, false, nil
			}
			if len(gotParams) != 0 {
				return functions.Shape{}, nil, nil, false, nil
			}
			if len(gotCall.Args) != 0 {
				return functions.Shape{}, nil, nil, false, nil
			}
			if gotCall.Name != "must" {
				return functions.Shape{}, nil, nil, false, nil
			}
			outputShape := functions.Shape{
				Type:  bindings["A"],
				Stack: gotInputShape.Stack,
			}
			action := func(inputState states.State, args []states.Action) *states.Thunk {
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
							errors.Pos(gotCall.Pos),
							errors.GotValue(inputState.Value)),
					)
				}
			}
			return outputShape, action, nil, true, nil
		},
	})
}
