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

func (x MappingExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return Shape{}, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// make sure the input type is a sequence type
	if !types.AnySeqType.Subsumes(inputShape.Type) {
		return Shape{}, nil, errors.E(
			errors.Code(errors.MappingRequiresSeqType),
			errors.Pos(x.Pos),
			errors.WantType(types.AnySeqType),
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
		Type:  &types.SeqType{bodyOutputShape.Type},
		Stack: inputShape.Stack,
	}
	// create action
	action := func(inputState states.State, args []states.Action) states.State {
		channel := make(chan values.Value)
		go func() {
			for el := range inputState.Value.Iter() {
				bodyInputState := states.State{
					Value: el,
					Stack: inputState.Stack,
				}
				bodyOutputState := bodyAction(bodyInputState, nil)
				channel <- bodyOutputState.Value
			}
			close(channel)
		}()
		return states.State{
			Value: values.SeqValue{bodyOutputShape.Type, channel},
			Stack: inputState.Stack,
		}
	}
	return outputShape, action, nil
}
