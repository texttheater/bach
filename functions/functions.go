package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/states"
)

type Expression interface {
	Position() lexer.Position
	Typecheck(inputShape Shape, params []*Parameter) (outputShape Shape, action states.Action, err error)
}
