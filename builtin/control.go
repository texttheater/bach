package builtin

import (
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initControl() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		func(gotInputShape functions.Shape, gotCall functions.CallExpression, gotParams []*functions.Parameter) (functions.Shape, states.Action, bool, error) {
			if len(gotParams) != 0 {
				return functions.Shape{}, nil, false, nil
			}
			if len(gotCall.Args) != 0 {
				return functions.Shape{}, nil, false, nil
			}
			if gotCall.Name != "fatal" {
				return functions.Shape{}, nil, false, nil
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
						errors.GotValue(inputState.Value),
					),
				)
			}
			return outputShape, action, true, nil
		},
		func(gotInputShape functions.Shape, gotCall functions.CallExpression, gotParams []*functions.Parameter) (functions.Shape, states.Action, bool, error) {
			if len(gotParams) != 0 {
				return functions.Shape{}, nil, false, nil
			}
			if len(gotCall.Args) != 0 {
				return functions.Shape{}, nil, false, nil
			}
			if gotCall.Name != "must" {
				return functions.Shape{}, nil, false, nil
			}
			u, ok := gotInputShape.Type.(types.UnionType)
			if !ok {
				return functions.Shape{}, nil, false, nil
			}
			nullTypeFound := false
			var outputType types.Type = types.VoidType{}
			for _, disjunct := range u {
				if (types.NullType{}).Subsumes(disjunct) {
					nullTypeFound = true
				} else {
					outputType = types.Union(outputType, disjunct)
				}
			}
			if !nullTypeFound {
				return functions.Shape{}, nil, false, nil
			}
			outputShape := functions.Shape{
				Type:  outputType,
				Stack: gotInputShape.Stack,
			}
			action := func(inputState states.State, args []states.Action) *states.Thunk {
				inhabits, err := inputState.Value.Inhabits(types.NullType{}, inputState.TypeStack)
				if err != nil {
					return states.ThunkFromError(err)
				}
				if inhabits {
					return states.ThunkFromError(
						errors.E(
							errors.Code(errors.UnexpectedValue),
							errors.Pos(gotCall.Pos),
							errors.GotType(gotInputShape.Type),
							errors.GotValue(inputState.Value),
						),
					)
				}
				return states.ThunkFromState(inputState)
			}
			return outputShape, action, true, nil
		},
	})
}
