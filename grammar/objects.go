package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/ast"
)

type Object struct {
	Pos    lexer.Position `"{"`
	Name   *string        `[ @Name`
	Value  *Composition   `  ":" @@`
	Names  []string       `  { "," @Name`
	Values []*Composition `    ":" @@ } ] "}"`
}

func (g *Object) Ast() ast.Expression {
	propValMap := make(map[string]ast.Expression)
	if g.Name != nil {
		propValMap[*g.Name] = g.Value.Ast()
		for i := range g.Names {
			propValMap[g.Names[i]] = g.Values[i].Ast()
		}
	}
	return &ast.ObjExpression{g.Pos, propValMap}
}
