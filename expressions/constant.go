package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type ConstantExpression struct {
	Pos   lexer.Position
	Type  types.Type
	Value states.Value
}

func (x ConstantExpression) Position() lexer.Position {
	return x.Pos
}

func (x ConstantExpression) Typecheck(inputShape shapes.Shape, params []*params.Param) (shapes.Shape, states.Action, *states.IDStack, error) {
	if len(params) > 0 {
		return shapes.Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	outputShape := shapes.Shape{
		Type:  x.Type,
		Stack: inputShape.Stack,
	}
	action := func(inputState states.State, args []states.Action) *states.Thunk {
		return states.ThunkFromState(states.State{
			Value:     x.Value,
			Stack:     inputState.Stack,
			TypeStack: inputState.TypeStack,
		})
	}
	return outputShape, action, nil, nil
}
