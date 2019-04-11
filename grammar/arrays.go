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

func (g *Array) Ast() (expressions.Expression, error) {
	var elements []expressions.Expression
	if g.Element != nil {
		elements = make([]expressions.Expression, 1+len(g.Elements))
		var err error
		elements[0], err = g.Element.Ast()
		if err != nil {
			return nil, err
		}
		for i, element := range g.Elements {
			elements[i+1], err = element.Ast()
			if err != nil {
				return nil, err
			}
		}
	}
	return &expressions.ArrExpression{g.Pos, elements}, nil
}
