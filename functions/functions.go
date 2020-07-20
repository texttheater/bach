package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/states"
)

type Expression interface {
	Position() lexer.Position
	Typecheck(inputShape Shape, params []*states.Parameter) (outputShape Shape, action states.Action, IDs *states.IDStack, err error)
}
