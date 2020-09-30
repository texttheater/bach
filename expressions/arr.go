package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type ArrExpression struct {
	Pos      lexer.Position
	Elements []Expression
	RestPos  lexer.Position
	Rest     Expression
}

func (x ArrExpression) Position() lexer.Position {
	return x.Pos
}

func (x ArrExpression) Typecheck(inputShape Shape, params []*parameters.Parameter) (Shape, states.Action, *states.IDStack, error) {
	// make sure we got no params
	if len(params) > 0 {
		return Shape{}, nil, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos))

	}
	// typecheck rest
	var outputType types.Type
	var action states.Action
	var ids *states.IDStack
	if x.Rest == nil {
		outputType = types.VoidArrType
		action = func(inputState states.State, args []states.Action) *states.Thunk {
			return states.ThunkFromState(states.State{
				Value: (*states.ArrValue)(nil),
				Stack: inputState.Stack,
			})
		}
	} else {
		var restShape Shape
		var restIDs *states.IDStack
		var err error
		restShape, action, restIDs, err = x.Rest.Typecheck(inputShape, nil)
		if err != nil {
			return Shape{}, nil, nil, err
		}
		if !(types.AnyArrType).Subsumes(restShape.Type) {
			return Shape{}, nil, nil, errors.E(
				errors.Code(errors.RestRequiresArrType),
				errors.Pos(x.RestPos),
				errors.WantType(types.AnyArrType),
				errors.GotType(restShape.Type))

		}
		ids = ids.AddAll(restIDs)
		outputType = restShape.Type
	}
	// typecheck elements
	for i := len(x.Elements) - 1; i >= 0; i-- {
		elementShape, elementAction, elementIDs, err := x.Elements[i].Typecheck(inputShape, nil)
		if err != nil {
			return Shape{}, nil, nil, err
		}
		outputType = &types.NearrType{
			HeadType: elementShape.Type,
			TailType: outputType,
		}
		tailAction := action
		action = func(inputState states.State, args []states.Action) *states.Thunk {
			res := elementAction(inputState, nil).Eval()
			if res.Error != nil {
				return states.ThunkFromError(res.Error)
			}
			return states.ThunkFromState(states.State{
				Value: &states.ArrValue{
					Head: res.Value,
					Tail: tailAction(inputState, nil),
				},
				Stack: inputState.Stack,
			})
		}
		ids = ids.AddAll(elementIDs)
	}
	// make output shape
	outputShape := Shape{
		Type:  outputType,
		Stack: inputShape.Stack,
	}
	// return
	return outputShape, action, ids, nil
}
