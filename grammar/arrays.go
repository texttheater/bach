package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
)

type Array struct {
	Pos      lexer.Position `"["`
	Element  *Composition   `( @@`
	Elements []*Composition `  ( "," @@ )* )? "]"`
}

func (g *Array) Ast() expressions.Expression {
	var elements []expressions.Expression
	if g.Element != nil {
		elements = make([]expressions.Expression, 1+len(g.Elements))
		elements[0] = g.Element.Ast()
		for i, element := range g.Elements {
			elements[i+1] = element.Ast()
		}
	}
	return &expressions.ArrExpression{g.Pos, elements}
}
