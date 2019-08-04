package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
)

type Parenthesis struct {
	Pos         lexer.Position
	Composition *Composition `"(" @@ ")"`
}

func (g *Parenthesis) Ast() (expressions.Expression, error) {
	return g.Composition.Ast()
}
