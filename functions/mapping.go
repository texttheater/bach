package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type MappingExpression struct {
	Pos  lexer.Position
	Body Expression
}

func (x MappingExpression) Position() lexer.Position {
	return x.Pos
}

func (x MappingExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return Shape{}, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// make sure the input type is a sequence type
	if !types.AnyArrType.Subsumes(inputShape.Type) {
		return Shape{}, nil, errors.E(
			errors.Code(errors.MappingRequiresArrType),
			errors.Pos(x.Pos),
			errors.WantType(types.AnyArrType),
			errors.GotType(inputShape.Type),
		)
	}
	// typecheck body
	bodyInputShape := Shape{
		Type:  inputShape.Type.ElementType(),
		Stack: inputShape.Stack,
	}
	bodyOutputShape, bodyAction, err := x.Body.Typecheck(bodyInputShape, nil)
	if err != nil {
		return Shape{}, nil, err
	}
	// create output shape
	outputShape := Shape{
		Type: &types.ArrType{
			ElType: bodyOutputShape.Type,
		},
		Stack: inputShape.Stack,
	}
	// create action
	action := func(inputState states.State, args []states.Action) *states.Thunk {
		arr := inputState.Value.(*states.ArrValue)
		var next func() *states.Thunk
		next = func() *states.Thunk {
			if arr == nil {
				return states.ThunkFromValue((*states.ArrValue)(nil))
			}
			bodyInputState := states.State{
				Value:     arr.Head,
				Stack:     inputState.Stack,
				TypeStack: inputState.TypeStack,
			}
			res := bodyAction(bodyInputState, nil).Eval()
			if res.Error != nil {
				return states.ThunkFromError(res.Error)
			}
			var err error
			arr, err = arr.GetTail()
			if err != nil {
				return states.ThunkFromError(err)
			}
			if res.Drop {
				return next()
			}
			return states.ThunkFromValue(
				&states.ArrValue{
					Head: res.State.Value,
					Tail: &states.Thunk{
						Func: func() *states.Thunk {
							return next()
						},
					},
				})

		}
		return next()
	}
	return outputShape, action, nil
}
