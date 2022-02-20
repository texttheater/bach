package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type Composition struct {
	Pos        lexer.Position
	Component  *Component   `@@`
	Components []*Component `( @@ )*`
}

func (g *Composition) Ast() (expressions.Expression, error) {
	pos := g.Component.Pos
	e, err := g.Component.Ast()
	if err != nil {
		return nil, err
	}
	for _, comp := range g.Components {
		compAst, err := comp.Ast()
		if err != nil {
			return nil, err
		}
		e = &expressions.CompositionExpression{pos, e, compAst}
	}
	return e, nil
}

type Component struct {
	Pos         lexer.Position
	NumLiteral  *float64     `  @NumLiteral`
	StrLiteral  *StrLiteral  `| @@`
	ArrLiteral  *ArrLiteral  `| @@`
	ObjLiteral  *ObjLiteral  `| @@`
	Call        *Call        `| @@`
	Assignment  *Assignment  `| @@`
	Definition  *Definition  `| @@`
	Conditional *Conditional `| @@`
	Regexp      *Regexp      `| @@`
	Getter      *Getter      `| @@`
}

func (g *Component) Ast() (expressions.Expression, error) {
	if g.NumLiteral != nil {
		return &expressions.ConstantExpression{
			Pos:   g.Pos,
			Type:  types.Num{},
			Value: states.NumValue(*g.NumLiteral),
		}, nil
	}
	if g.StrLiteral != nil {
		return g.StrLiteral.Ast()
	}
	if g.ArrLiteral != nil {
		return g.ArrLiteral.Ast()
	}
	if g.ObjLiteral != nil {
		return g.ObjLiteral.Ast()
	}
	if g.Call != nil {
		return g.Call.Ast()
	}
	if g.Assignment != nil {
		return g.Assignment.Ast()
	}
	if g.Definition != nil {
		return g.Definition.Ast()
	}
	if g.Conditional != nil {
		return g.Conditional.Ast()
	}
	if g.Regexp != nil {
		return g.Regexp.Ast()
	}
	if g.Getter != nil {
		return g.Getter.Ast()
	}
	panic("invalid component")
}
