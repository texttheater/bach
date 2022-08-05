package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
)

type Assignment struct {
	Pos    lexer.Position
	EqName *EqName `  @@`
}

func (g *Assignment) Ast() (expressions.Expression, error) {
	if g.EqName != nil {
		return g.EqName.Ast()
	}
	panic("invalid assignment")
}

type EqName struct {
	Pos    lexer.Position
	EqName string `@EqName`
}

func (g *EqName) Ast() (expressions.Expression, error) {
	name := g.EqName[1:]
	return &expressions.AssignmentExpression{
		Pos:  g.Pos,
		Name: name,
	}, nil
}
