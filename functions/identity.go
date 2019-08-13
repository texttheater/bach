package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/states"
)

type IdentityExpression struct {
	Pos lexer.Position
}

func (x IdentityExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	return inputShape, func(inputState states.State, args []states.Action) states.State {
		return inputState
	}, nil
}
