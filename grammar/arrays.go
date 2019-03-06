package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/ast"
)

type Array struct {
	Pos      lexer.Position `"["`
	Element  *Composition   `[ @@`
	Elements []*Composition `  { "," @@ } ] "]"`
}

func (g *Array) Ast() ast.Expression {
	var elements []ast.Expression
	if g.Element != nil {
		elements = make([]ast.Expression, 1+len(g.Elements))
		elements[0] = g.Element.Ast()
		for i, element := range g.Elements {
			elements[i+1] = element.Ast()
		}
	}
	return &ast.ArrExpression{g.Pos, elements}
}
