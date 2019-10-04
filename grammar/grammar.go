package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type Composition struct {
	Pos        lexer.Position
	Component  *SComponent   `@@`
	Components []*SComponent `( @@ )*`
}

func (g *Composition) Ast() (functions.Expression, error) {
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
		e = &functions.CompositionExpression{pos, e, compAst}
	}
	return e, nil
}

type SComponent struct {
	Pos         lexer.Position
	Num         *float64     `  @Num`
	Str         *string      `| @Str`
	Array       *Array       `| @@`
	Object      *Object      `| @@`
	Call        *Call        `| @@`
	Assignment  *Assignment  `| @Assignment`
	Definition  *Definition  `| @@`
	Conditional *Conditional `| @@`
	Drop        *Drop        `| @@`
	Filter      *Filter      `| @@`
	Regexp      *Regexp      `| @@`
}

func (g *SComponent) Ast() (functions.Expression, error) {
	if g.Num != nil {
		return &functions.ConstantExpression{
			Pos:   g.Pos,
			Type:  types.NumType{},
			Value: values.NumValue(*g.Num),
		}, nil
	}
	if g.Str != nil {
		return &functions.ConstantExpression{
			Pos:   g.Pos,
			Type:  types.StrType{},
			Value: values.StrValue(*g.Str),
		}, nil
	}
	if g.Array != nil {
		return g.Array.Ast()
	}
	if g.Object != nil {
		return g.Object.Ast()
	}
	if g.Call != nil {
		return g.Call.Ast()
	}
	if g.Regexp != nil {
		return g.Regexp.Ast()
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
	if g.Drop != nil {
		return g.Drop.Ast()
	}
	if g.Filter != nil {
		return g.Filter.Ast()
	}
	panic("invalid component")
}

type PComponent struct {
	Pos        lexer.Position
	Num        *float64    `  @Num`
	Str        *string     `| @Str`
	Array      *Array      `| @@`
	Object     *Object     `| @@`
	Call       *Call       `| @@`
	Assignment *Assignment `| @Assignment`
	Definition *Definition `| @@`
	Drop       *Drop       `| @@`
	Filter     *Filter     `| @@`
	Regexp     *Regexp     `| @@`
}

func (g *PComponent) Ast() (functions.Expression, error) {
	if g.Num != nil {
		return &functions.ConstantExpression{
			Pos:   g.Pos,
			Type:  types.NumType{},
			Value: values.NumValue(*g.Num),
		}, nil
	}
	if g.Str != nil {
		return &functions.ConstantExpression{
			Pos:   g.Pos,
			Type:  types.StrType{},
			Value: values.StrValue(*g.Str),
		}, nil
	}
	if g.Array != nil {
		return g.Array.Ast()
	}
	if g.Object != nil {
		return g.Object.Ast()
	}
	if g.Call != nil {
		return g.Call.Ast()
	}
	if g.Regexp != nil {
		return g.Regexp.Ast()
	}
	if g.Assignment != nil {
		return g.Assignment.Ast()
	}
	if g.Definition != nil {
		return g.Definition.Ast()
	}
	if g.Drop != nil {
		return g.Drop.Ast()
	}
	if g.Filter != nil {
		return g.Filter.Ast()
	}
	panic("invalid component")
}
