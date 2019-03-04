package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/ast"
)

type Object struct {
	Pos    lexer.Position `"{"`
	Prop   *string        `[ @Prop`
	Value  *Composition   `  ":" @@`
	Props  []string       `  { "," @Prop`
	Values []*Composition `    ":" @@ } ] "}"`
}

func (g *Object) Ast() ast.Expression {
	propValMap := make(map[string]ast.Expression)
	if g.Prop != nil {
		propValMap[*g.Prop] = g.Value.Ast()
		for i := range g.Props {
			propValMap[g.Props[i]] = g.Values[i].Ast()
		}
	}
	return &ast.ObjExpression{g.Pos, propValMap}
}
