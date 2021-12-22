package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

// A composition consists of components. There are different types of
// components which can appear in different contexts:
// S-components appear at the top level.
// P-components appear within mappings. Compared to S-components, they
// exclude conditionals (these are tied in with mapping syntax) and add
// Drop (within a mapping, you can drop elements).
// Q-components appear within branches within mappings. Compared to
// P-components, they add conditionals (we are now one level deeper and
// therefore conditionals are fine again) but retain Drop (we are still in
// a mapping).

type Composition struct {
	Pos        lexer.Position
	Component  *SComponent   `@@`
	Components []*SComponent `( @@ )*`
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

type SComponent struct {
	Pos         lexer.Position
	NumLiteral  *float64     `  @NumLiteral`
	StrLiteral  *StrLiteral  `| @@`
	ArrLiteral  *ArrLiteral  `| @@`
	ObjLiteral  *ObjLiteral  `| @@`
	Call        *Call        `| @@`
	Assignment  *Assignment  `| @@`
	Definition  *Definition  `| @@`
	Conditional *Conditional `| @@`
	Filter      *Filter      `| @@`
	Regexp      *Regexp      `| @@`
	Getter      *Getter      `| @@`
}

func (g *SComponent) Ast() (expressions.Expression, error) {
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
	NumLiteral *float64    `  @NumLiteral`
	StrLiteral *StrLiteral `| @@`
	ArrLiteral *ArrLiteral `| @@`
	ObjLiteral *ObjLiteral `| @@`
	Call       *Call       `| @@`
	Assignment *Assignment `| @@`
	Definition *Definition `| @@`
	Filter     *Filter     `| @@`
	Regexp     *Regexp     `| @@`
	Getter     *Getter     `| @@`
	Drop       *Drop       `| @@`
}

func (g *PComponent) Ast() (expressions.Expression, error) {
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

func (g *QComposition) Ast() (expressions.Expression, error) {
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

type QComponent struct {
	Pos         lexer.Position
	NumLiteral  *float64     `  @NumLiteral`
	StrLiteral  *StrLiteral  `| @@`
	ArrLiteral  *ArrLiteral  `| @@`
	ObjLiteral  *ObjLiteral  `| @@`
	Call        *Call        `| @@`
	Assignment  *Assignment  `| @@`
	Definition  *Definition  `| @@`
	Conditional *Conditional `| @@`
	Filter      *Filter      `| @@`
	Regexp      *Regexp      `| @@`
	Getter      *Getter      `| @@`
	Drop        *Drop        `| @@`
}

func (g *QComponent) Ast() (expressions.Expression, error) {
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
