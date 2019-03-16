package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
)

type Mapping struct {
	Pos  lexer.Position
	Body *Composition `"each" @@ "all"`
}

func (g *Mapping) Ast() expressions.Expression {
	return &expressions.MappingExpression{
		Pos:  g.Pos,
		Body: g.Body.Ast(),
	}
}
