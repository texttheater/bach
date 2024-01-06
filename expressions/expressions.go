package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
)

type Expression interface {
	Position() lexer.Position
	Typecheck(inputShape shapes.Shape, params []*params.Param) (outputShape shapes.Shape, action states.Action, IDs *states.IDStack, err error)
}
