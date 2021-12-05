package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
)

type ArrLiteral struct {
	Pos  lexer.Position  `"["`
	Rest *ArrLiteralRest "@@"
}

func (g *ArrLiteral) Ast() (expressions.Expression, error) {
	ast, err := g.Rest.Ast()
	if ast != nil {
		ast.Pos = g.Pos
	}
	return ast, err
}

type ArrLiteralRest struct {
	Element  *Composition   `( @@`
	Elements []*Composition `  ( "," @@ )*`
	Rest     *Composition   `  ( ";" @@ )? )? "]"`
}

func (g *ArrLiteralRest) Ast() (*expressions.ArrExpression, error) {
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
		Elements: elements,
		Rest:     rest,
	}
	if g.Rest != nil {
		x.RestPos = g.Rest.Pos
	}
	return x, nil
}
