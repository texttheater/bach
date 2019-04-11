package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
)

type Object struct {
	Pos    lexer.Position `"{"`
	Prop   *string        `( @Prop`
	Value  *Composition   `  ":" @@`
	Props  []string       `  ( "," @Prop`
	Values []*Composition `    ":" @@ )* )? "}"`
}

func (g *Object) Ast() (expressions.Expression, error) {
	propValMap := make(map[string]expressions.Expression)
	if g.Prop != nil {
		var err error
		propValMap[*g.Prop], err = g.Value.Ast()
		if err != nil {
			return nil, err
		}
		for i := range g.Props {
			propValMap[g.Props[i]], err = g.Values[i].Ast()
			if err != nil {
				return nil, err
			}
		}
	}
	return &expressions.ObjExpression{g.Pos, propValMap}, nil
}
