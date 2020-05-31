package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

// A composition consists of components. There are different types of
// components which can appear in different contexts:
// S-components appear at the top level.
// P-components appear within mappings. Compared to S-components, they
// exclude conditionals (these are tied in with mapping syntax) and add
// Drop (within a mapping, you can drop elements).
// Q-componnts appear within branches within mappings. Compared to
// P-components, they add conditionals (we are now one level deeper and
// therefore conditionals are fine again) but retain Drop (we are still in
// a mapping).

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
	Assignment  *Assignment  `| @@`
	Definition  *Definition  `| @@`
	Conditional *Conditional `| @@`
	Filter      *Filter      `| @@`
	Regexp      *Regexp      `| @@`
	Getter      *Getter      `| @@`
}

func (g *SComponent) Ast() (functions.Expression, error) {
	if g.Num != nil {
		return &functions.ConstantExpression{
			Pos:   g.Pos,
			Type:  types.NumType{},
			Value: states.NumValue(*g.Num),
		}, nil
	}
	if g.Str != nil {
		return &functions.ConstantExpression{
			Pos:   g.Pos,
			Type:  types.StrType{},
			Value: states.StrValue(*g.Str),
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
	if g.Assignment != nil {
		return g.Assignment.Ast()
	}
	if g.Definition != nil {
		return g.Definition.Ast()
	}
	if g.Conditional != nil {
		return g.Conditional.Ast()
	}
	if g.Filter != nil {
		return g.Filter.Ast()
	}
	if g.Regexp != nil {
		return g.Regexp.Ast()
	}
	if g.Getter != nil {
		return g.Getter.Ast()
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
	Assignment *Assignment `| @@`
	Definition *Definition `| @@`
	Filter     *Filter     `| @@`
	Regexp     *Regexp     `| @@`
	Getter     *Getter     `| @@`
	Drop       *Drop       `| @@`
}

func (g *PComponent) Ast() (functions.Expression, error) {
	if g.Num != nil {
		return &functions.ConstantExpression{
			Pos:   g.Pos,
			Type:  types.NumType{},
			Value: states.NumValue(*g.Num),
		}, nil
	}
	if g.Str != nil {
		return &functions.ConstantExpression{
			Pos:   g.Pos,
			Type:  types.StrType{},
			Value: states.StrValue(*g.Str),
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
	if g.Assignment != nil {
		return g.Assignment.Ast()
	}
	if g.Definition != nil {
		return g.Definition.Ast()
	}
	if g.Filter != nil {
		return g.Filter.Ast()
	}
	if g.Regexp != nil {
		return g.Regexp.Ast()
	}
	if g.Getter != nil {
		return g.Getter.Ast()
	}
	if g.Drop != nil {
		return g.Drop.Ast()
	}
	panic("invalid component")
}

type QComposition struct {
	Pos        lexer.Position
	Component  *QComponent   `@@`
	Components []*QComponent `( @@ )*`
}

func (g *QComposition) Ast() (functions.Expression, error) {
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

type QComponent struct {
	Pos         lexer.Position
	Num         *float64     `  @Num`
	Str         *string      `| @Str`
	Array       *Array       `| @@`
	Object      *Object      `| @@`
	Call        *Call        `| @@`
	Assignment  *Assignment  `| @@`
	Definition  *Definition  `| @@`
	Conditional *Conditional `| @@`
	Filter      *Filter      `| @@`
	Regexp      *Regexp      `| @@`
	Getter      *Getter      `| @@`
	Drop        *Drop        `| @@`
}

func (g *QComponent) Ast() (functions.Expression, error) {
	if g.Num != nil {
		return &functions.ConstantExpression{
			Pos:   g.Pos,
			Type:  types.NumType{},
			Value: states.NumValue(*g.Num),
		}, nil
	}
	if g.Str != nil {
		return &functions.ConstantExpression{
			Pos:   g.Pos,
			Type:  types.StrType{},
			Value: states.StrValue(*g.Str),
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
	if g.Assignment != nil {
		return g.Assignment.Ast()
	}
	if g.Definition != nil {
		return g.Definition.Ast()
	}
	if g.Conditional != nil {
		return g.Conditional.Ast()
	}
	if g.Filter != nil {
		return g.Filter.Ast()
	}
	if g.Regexp != nil {
		return g.Regexp.Ast()
	}
	if g.Getter != nil {
		return g.Getter.Ast()
	}
	if g.Drop != nil {
		return g.Drop.Ast()
	}
	panic("invalid component")
}
