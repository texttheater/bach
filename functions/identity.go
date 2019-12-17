package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/states"
)

type IdentityExpression struct {
	Pos lexer.Position
}

func (x IdentityExpression) Position() lexer.Position {
	return x.Pos
}

func (x IdentityExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, *states.IDStack, error) {
	return inputShape, func(inputState states.State, args []states.Action) *states.Thunk {
		return states.ThunkFromState(inputState)
	}, nil, nil
}
