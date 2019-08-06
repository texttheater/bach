package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

type Pattern struct {
	Pos         lexer.Position
	NamePattern *NamePattern `  @@`
	TypePattern *TypePattern `| @@`
	ArrPattern  *ArrPattern  `| @@`
	ObjPattern  *ObjPattern  `| @@`
}

func (g *Pattern) Ast() (functions.Pattern, error) {
	if g.NamePattern != nil {
		p, err := g.NamePattern.Ast()
		if err != nil {
			return nil, err
		}
		return p, nil
	} else if g.TypePattern != nil {
		p, err := g.TypePattern.Ast()
		if err != nil {
			return nil, err
		}
		return p, nil
	} else if g.ArrPattern != nil {
		p, err := g.ArrPattern.Ast()
		if err != nil {
			return nil, err
		}
		return p, nil
	} else if g.ObjPattern != nil {
		p, err := g.ObjPattern.Ast()
		if err != nil {
			return nil, err
		}
		return p, nil
	} else {
		panic("invalid pattern")
	}
}

type NamePattern struct {
	Pos  lexer.Position
	Name *string `@Prop | @Op1 | @Op2`
}

func (g *NamePattern) Ast() (functions.Pattern, error) {
	return functions.TypePattern{g.Pos, types.AnyType{}, g.Name}, nil
}

type TypePattern struct {
	Pos  lexer.Position
	Type *Type   `@@`
	Name *string `( @Prop | @Op1 | @Op2 )?`
}

func (g *TypePattern) Ast() (functions.Pattern, error) {
	return functions.TypePattern{g.Pos, g.Type.Ast(), g.Name}, nil
}

type ArrPattern struct {
	Pos      lexer.Position `"["`
	Element  *Pattern       `( @@`
	Elements []*Pattern     `  ( "," @@ )* )? "]"`
	Name     *string        `( @Prop | @Op1 | @Op2 )?`
}

func (g *ArrPattern) Ast() (functions.Pattern, error) {
	var elPatterns []functions.Pattern
	if g.Element != nil {
		elPatterns = make([]functions.Pattern, len(g.Elements)+1)
		p, err := g.Element.Ast()
		if err != nil {
			return nil, err
		}
		elPatterns[0] = p
		for i, el := range g.Elements {
			p, err = el.Ast()
			if err != nil {
				return nil, err
			}
			elPatterns[i+1] = p
		}
	}
	return &functions.ArrPattern{g.Pos, elPatterns, g.Name}, nil
}

type ObjPattern struct {
	Pos    lexer.Position `"{"`
	Prop   *string        `( @Prop`
	Value  *Pattern       `  ":" @@`
	Props  []string       `   ( "," @Prop`
	Values []*Pattern     `     ":" @@ )* )? "}"`
	Name   *string        `( @Prop | @Op1 | @Op2 )?`
}

func (g *ObjPattern) Ast() (functions.Pattern, error) {
	propPatternMap := make(map[string]functions.Pattern)
	if g.Prop != nil {
		p, err := g.Value.Ast()
		if err != nil {
			return nil, err
		}
		propPatternMap[*g.Prop] = p
		for i, prop := range g.Props {
			p, err := g.Values[i].Ast()
			if err != nil {
				return nil, err
			}
			propPatternMap[prop] = p
		}
	}
	return &functions.ObjPattern{g.Pos, propPatternMap, g.Name}, nil
}
