package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
)

type Array struct {
	Pos      lexer.Position `"["`
	Element  *Composition   `( @@`
	Elements []*Composition `  ( "," @@ )*`
	Rest     *Composition   `  ( ";" @@ )? )? "]"`
}

func (g *Array) Ast() (expressions.Expression, error) {
	var elements []expressions.Expression
	var rest expressions.Expression
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
		if g.Rest != nil {
			rest, err = g.Rest.Ast()
			if err != nil {
				return nil, err
			}
		}
	}
	x := &expressions.ArrExpression{
		Pos:      g.Pos,
		Elements: elements,
		Rest:     rest,
	}
	if g.Rest != nil {
		x.RestPos = g.Rest.Pos
	}
	return x, nil
}
