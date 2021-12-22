package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/states"
)

type Expression interface {
	Position() lexer.Position
	Typecheck(inputShape Shape, params []*params.Param) (outputShape Shape, action states.Action, IDs *states.IDStack, err error)
}
