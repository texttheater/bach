package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/expressions"
)

type ObjLiteral struct {
	Pos    lexer.Position `"{"`
	Prop   *Prop          `( @@`
	Value  *Composition   `  ":" @@`
	Props  []*Prop        `  ( "," @@`
	Values []*Composition `    ":" @@ )* )? "}"`
}

func (g *ObjLiteral) Ast() (expressions.Expression, error) {
	propValMap := make(map[string]expressions.Expression)
	if g.Prop != nil {
		var err error
		prop, err := g.Prop.StaticStr()
		if err != nil {
			return nil, err
		}
		propValMap[prop], err = g.Value.Ast()
		if err != nil {
			return nil, err
		}
		for i := range g.Props {
			prop, err = g.Props[i].StaticStr()
			if err != nil {
				return nil, err
			}
			propValMap[prop], err = g.Values[i].Ast()
			if err != nil {
				return nil, err
			}
		}
	}
	return &expressions.ObjExpression{g.Pos, propValMap}, nil
}

type Prop struct {
	Pos        lexer.Position
	StrLiteral *StrLiteral `  @@`
	Other      *string     `| @Lid | @Op1 | @Op2 | @NumLiteral`
}

func (g *Prop) StaticStr() (string, error) {
	if g.StrLiteral != nil {
		str, ok, err := g.StrLiteral.StaticStr()
		if err != nil {
			return "", err
		}
		if !ok {
			return "", errors.E(
				errors.Code(errors.Syntax),
				errors.Pos(g.Pos),
				errors.Message("Can't use a dynamic string literal for object property."),
			)
		}
		return str, nil
	}
	return *g.Other, nil
}
