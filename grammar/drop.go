package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
)

type Drop struct {
	Pos lexer.Position `"drop"`
}

func (g *Drop) Ast() (expressions.Expression, error) {
	return &expressions.DropExpression{
		Pos: g.Pos,
	}, nil
}
