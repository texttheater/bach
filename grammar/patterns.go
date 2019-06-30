package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/patterns"
)

type Pattern struct {
	Pos         lexer.Position
	TypePattern *TypePattern `  @@`
	ArrPattern  *ArrPattern  `| @@`
	ObjPattern  *ObjPattern  `| @@`
}

func (g *Pattern) Ast() (patterns.Pattern, error) {
	if g.TypePattern != nil {
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

type TypePattern struct {
	Pos  lexer.Position
	Type *Type   `@@`
	Name *string `( @Prop | @Op1 | @Op2 )?`
}

func (g *TypePattern) Ast() (patterns.Pattern, error) {
	return patterns.TypePattern{g.Pos, g.Type.Ast(), g.Name}, nil
}

type ArrPattern struct {
	Pos      lexer.Position `"["`
	Element  *Pattern       `( @@`
	Elements []*Pattern     `  ( "," @@ )* )? "]"`
	Name     *string        `( @Prop | @Op1 | @Op2 )?`
}

func (g *ArrPattern) Ast() (patterns.Pattern, error) {
	var elPatterns []patterns.Pattern
	if g.Element != nil {
		elPatterns = make([]patterns.Pattern, len(g.Elements)+1)
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
	return &patterns.ArrPattern{g.Pos, elPatterns, g.Name}, nil
}

type ObjPattern struct {
	Pos    lexer.Position `"{"`
	Prop   *string        `( @Prop`
	Value  *Pattern       `  ":" @@`
	Props  []string       `   ( "," @Prop`
	Values []*Pattern     `     ":" @@ )* )? "}"`
	Name   *string        `( @Prop | @Op1 | @Op2 )?`
}

func (g *ObjPattern) Ast() (patterns.Pattern, error) {
	propPatternMap := make(map[string]patterns.Pattern)
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
	return &patterns.ObjPattern{g.Pos, propPatternMap, g.Name}, nil
}
