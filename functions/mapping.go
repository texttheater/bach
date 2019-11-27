package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
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
	action := func(inputState states.State, args []states.Action) states.Thunk {
		arr := inputState.Value.(*values.ArrValue)
		var next func() (values.Value, *values.ArrValue, error)
		next = func() (values.Value, *values.ArrValue, error) {
			err := arr.Eval()
			if err != nil {
				return nil, nil, err
			}
			if arr.Head == nil {
				return nil, nil, nil
			}
			bodyInputState := states.State{
				Value:     arr.Head,
				Stack:     inputState.Stack,
				TypeStack: inputState.TypeStack,
			}
			bodyOutputState, drop, err := bodyAction(bodyInputState, nil).Eval()
			if err != nil {
				return nil, nil, err
			}
			arr = arr.Tail
			if drop {
				return next()
			}
			return bodyOutputState.Value, &values.ArrValue{
				Func: next,
			}, nil
		}
		return states.Thunk{State: states.State{
			Value: &values.ArrValue{
				Func: next,
			},
			Stack:     inputState.Stack,
			TypeStack: inputState.TypeStack,
		}, Drop: false, Err: nil}

	}
	return outputShape, action, nil
}
