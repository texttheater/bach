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
		elements = make([]ast.Expression, 0, 1+len(g.Elements))
		elements = append(elements, g.Element.Ast())
		for _, element := range g.Elements {
			elements = append(elements, element.Ast())
		}
	}
	return &ast.ArrayExpression{g.Pos, elements}
}
