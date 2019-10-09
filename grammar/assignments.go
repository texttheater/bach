package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
)

type Assignment struct {
	Pos        lexer.Position
	Assignment string `@Assignment`
}

func (g *Assignment) Ast() (functions.Expression, error) {
	name := g.Assignment[1:]
	return &functions.AssignmentExpression{
		Pos:  g.Pos,
		Name: name,
	}, nil
}
