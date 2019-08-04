package expressions

import (
	"regexp/syntax"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
)

type RegexpExpression struct {
	Pos    lexer.Position
	Regexp *syntax.Regexp
}

func (x RegexpExpression) Typecheck(inputShape shapes.Shape, params []*shapes.Parameter) (shapes.Shape, states.Action, error) {
	panic("not implemented yet")
}
