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
	action := func(inputState states.State, args []states.Action) states.State {
		channel := make(chan values.Value)
		go func() {
			for el := range inputState.Value.Iter() {
				bodyInputState := states.State{
					Value:     el,
					Stack:     inputState.Stack,
					TypeStack: inputState.TypeStack,
				}
				bodyOutputState := bodyAction(bodyInputState, nil)
				if !bodyOutputState.Drop {
					channel <- bodyOutputState.Value
				}
			}
			close(channel)
		}()
		return states.State{
			Value: &values.ArrValue{
				Channel: channel,
			},
			Stack:     inputState.Stack,
			TypeStack: inputState.TypeStack,
		}
	}
	return outputShape, action, nil
}
