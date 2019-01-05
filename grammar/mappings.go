package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/ast"
)

type Mapping struct {
	Pos  lexer.Position
	Body *Composition `"each" @@ "all"`
}

func (g *Mapping) Ast() ast.Expression {
	return &ast.MappingExpression{
		Pos:  g.Pos,
		Body: g.Body.Ast(),
	}
}
