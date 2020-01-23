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

func (x MappingExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, *states.IDStack, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return Shape{}, nil, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// make sure the input type is a sequence type
	if !types.AnyArrType.Subsumes(inputShape.Type) {
		return Shape{}, nil, nil, errors.E(
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
	bodyOutputShape, bodyAction, bodyIDs, err := x.Body.Typecheck(bodyInputShape, nil)
	if err != nil {
		return Shape{}, nil, nil, err
	}
	ids := bodyIDs
	// create output shape
	outputShape := Shape{
		Type: &types.ArrType{
			ElType: bodyOutputShape.Type,
		},
		Stack: inputShape.Stack,
	}
	// create action
	action := func(inputState states.State, args []states.Action) *states.Thunk {
		input := states.IterFromValue(inputState.Value)
		var output func() (states.Value, bool, error)
		output = func() (states.Value, bool, error) {
			inputValue, ok, err := input()
			if err != nil {
				return nil, false, err
			}
			if !ok {
				return nil, false, nil
			}
			bodyInputState := states.State{
				Value:     inputValue,
				Stack:     inputState.Stack,
				TypeStack: inputState.TypeStack,
			}
			res := bodyAction(bodyInputState, nil).Eval()
			if res.Error != nil {
				return nil, false, res.Error
			}
			if res.Drop {
				return output()
			}
			return res.Value, true, nil
		}
		return states.ThunkFromIter(output)
	}
	return outputShape, action, ids, nil
}
