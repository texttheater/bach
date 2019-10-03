package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
)

type Drop struct {
	Pos lexer.Position `"drop"`
}

func (g *Drop) Ast() (functions.Expression, error) {
	return &functions.DropExpression{
		Pos: g.Pos,
	}, nil
}
