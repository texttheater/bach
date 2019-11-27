package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
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

func (x ArrExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	// make sure we got no params
	if len(params) > 0 {
		return Shape{}, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// typecheck rest
	var outputType types.Type
	var action states.Action
	if x.Rest == nil {
		outputType = types.VoidArrType
		action = func(inputState states.State, args []states.Action) states.Thunk {
			return states.Thunk{State: states.State{
				Value: &values.ArrValue{},
				Stack: inputState.Stack,
			}, Drop: false, Err: nil}

		}
	} else {
		var restShape Shape
		var err error
		restShape, action, err = x.Rest.Typecheck(inputShape, nil)
		if err != nil {
			return Shape{}, nil, err
		}
		if !(types.AnyArrType).Subsumes(restShape.Type) {
			return Shape{}, nil, errors.E(
				errors.Code(errors.RestRequiresArrType),
				errors.Pos(x.RestPos),
				errors.WantType(types.AnyArrType),
				errors.GotType(restShape.Type),
			)
		}
		outputType = restShape.Type
	}
	// typecheck elements
	for i := len(x.Elements) - 1; i >= 0; i-- {
		elementShape, elementAction, err := x.Elements[i].Typecheck(inputShape, nil)
		if err != nil {
			return Shape{}, nil, err
		}
		outputType = &types.NearrType{
			HeadType: elementShape.Type,
			TailType: outputType,
		}
		tailAction := action
		action = func(inputState states.State, args []states.Action) states.Thunk {
			return states.Thunk{State: states.State{
				Value: &values.ArrValue{
					Func: func() (values.Value, *values.ArrValue, error) {
						headState, _, err := elementAction(inputState, nil).Eval()
						if err != nil {
							return nil, nil, err
						}
						tailState, _, err := tailAction(inputState, nil).Eval()
						if err != nil {
							return nil, nil, err
						}
						return headState.Value, tailState.Value.(*values.ArrValue), nil
					},
				},
				Stack: inputState.Stack,
			}, Drop: false, Err: nil}

		}
	}
	// make output shape
	outputShape := Shape{
		Type:  outputType,
		Stack: inputShape.Stack,
	}
	// return
	return outputShape, action, nil
}
