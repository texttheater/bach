package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/ast"
)

///////////////////////////////////////////////////////////////////////////////

type Assignment struct {
	Pos  lexer.Position
	Name string
}

func (g *Assignment) Capture(values []string) error {
	g.Name = values[0][1:]
	return nil
}

func (g *Assignment) Ast() ast.Expression {
	return &ast.AssignmentExpression{
		Pos:  g.Pos,
		Name: g.Name,
	}
}

///////////////////////////////////////////////////////////////////////////////