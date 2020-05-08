package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
)

type Object struct {
	Pos    lexer.Position `"{"`
	Prop   *string        `( ( @Lid | @Op1 | @Op2 | @Num | @Str )`
	Value  *Composition   `  ":" @@`
	Props  []string       `  ( "," ( @Lid | @Op1 | @Op2 | @Num | @Str )`
	Values []*Composition `    ":" @@ )* )? "}"`
}

func (g *Object) Ast() (functions.Expression, error) {
	propValMap := make(map[string]functions.Expression)
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
	return &functions.ObjExpression{g.Pos, propValMap}, nil
}
