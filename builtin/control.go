package builtin

import (
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initControl() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		func(gotInputShape functions.Shape, gotCall functions.CallExpression, gotParams []*functions.Parameter) (functions.Shape, states.Action, *states.IDStack, bool, error) {
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
						errors.GotValue(inputState.Value),
					),
				)
			}
			return outputShape, action, nil, true, nil
		},
		func(gotInputShape functions.Shape, gotCall functions.CallExpression, gotParams []*functions.Parameter) (functions.Shape, states.Action, *states.IDStack, bool, error) {
			if len(gotParams) != 0 {
				return functions.Shape{}, nil, nil, false, nil
			}
			if len(gotCall.Args) != 0 {
				return functions.Shape{}, nil, nil, false, nil
			}
			if gotCall.Name != "must" {
				return functions.Shape{}, nil, nil, false, nil
			}
			u, ok := gotInputShape.Type.(types.UnionType)
			if !ok {
				return functions.Shape{}, nil, nil, false, nil
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
				return functions.Shape{}, nil, nil, false, nil
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
			return outputShape, action, nil, true, nil
		},
		func(gotInputShape functions.Shape, gotCall functions.CallExpression, gotParams []*functions.Parameter) (functions.Shape, states.Action, *states.IDStack, bool, error) {
			if len(gotParams) != 0 {
				return functions.Shape{}, nil, nil, false, nil
			}
			if len(gotCall.Args) != 0 {
				return functions.Shape{}, nil, nil, false, nil
			}
			if gotCall.Name != "must" {
				return functions.Shape{}, nil, nil, false, nil
			}
			u, ok := gotInputShape.Type.(types.UnionType)
			if !ok {
				return functions.Shape{}, nil, nil, false, nil
			}
			if len(u) != 2 {
				return functions.Shape{}, nil, nil, false, nil
			}
			left, ok := u[0].(types.ObjType)
			if !ok {
				return functions.Shape{}, nil, nil, false, nil
			}
			_, ok = left.PropTypeMap["left"]
			if !ok {
				return functions.Shape{}, nil, nil, false, nil
			}
			right, ok := u[1].(types.ObjType)
			if !ok {
				return functions.Shape{}, nil, nil, false, nil
			}
			if !ok {
				return functions.Shape{}, nil, nil, false, nil
			}
			rightType, ok := right.PropTypeMap["right"]
			outputShape := functions.Shape{
				Type:  rightType,
				Stack: gotInputShape.Stack,
			}
			action := func(inputState states.State, args []states.Action) *states.Thunk {
				inhabits, err := inputState.Value.Inhabits(right, inputState.TypeStack)
				if err != nil {
					return states.ThunkFromError(err)
				}
				if !inhabits {
					return states.ThunkFromError(
						errors.E(
							errors.Code(errors.UnexpectedValue),
							errors.Pos(gotCall.Pos),
							errors.GotType(gotInputShape.Type),
							errors.GotValue(inputState.Value),
						),
					)
				}
				return inputState.Value.(states.ObjValue)["right"]
			}
			return outputShape, action, nil, true, nil
		},
	})
}
