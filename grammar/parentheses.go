package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
)

type Parenthesis struct {
	Pos         lexer.Position
	Composition *Composition `"(" @@ ")"`
}

func (g *Parenthesis) Ast() (functions.Expression, error) {
	return g.Composition.Ast()
}
