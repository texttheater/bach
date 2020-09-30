package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
)

type Assignment struct {
	Pos        lexer.Position
	Assignment string `@Assignment`
}

func (g *Assignment) Ast() (expressions.Expression, error) {
	name := g.Assignment[1:]
	return &expressions.AssignmentExpression{
		Pos:  g.Pos,
		Name: name,
	}, nil
}
