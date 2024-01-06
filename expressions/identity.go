package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
)

type IdentityExpression struct {
	Pos lexer.Position
}

func (x IdentityExpression) Position() lexer.Position {
	return x.Pos
}

func (x IdentityExpression) Typecheck(inputShape shapes.Shape, params []*params.Param) (shapes.Shape, states.Action, *states.IDStack, error) {
	if len(params) > 0 {
		return shapes.Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	return inputShape, func(inputState states.State, args []states.Action) *states.Thunk {
		return states.ThunkFromState(inputState)
	}, nil, nil
}
