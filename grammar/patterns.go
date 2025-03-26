package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/types"
)

type Pattern struct {
	Pos         lexer.Position
	NamePattern *NamePattern `  @@`
	TypePattern *TypePattern `| @@`
	ArrPattern  *ArrPattern  `| @@`
	ObjPattern  *ObjPattern  `| @@`
}

func (g *Pattern) Ast() (expressions.Pattern, error) {
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
	Name string `@Lid | @Op1 | @Op2`
}

func (g *NamePattern) Ast() (expressions.Pattern, error) {
	return expressions.TypePattern{g.Pos, types.Any{}, &g.Name}, nil
}

type TypePattern struct {
	Pos  lexer.Position
	Type *types.TypeSyntax `@@`
	Name *string           `( @Lid | @Op1 | @Op2 )?`
}

func (g *TypePattern) Ast() (expressions.Pattern, error) {
	return expressions.TypePattern{g.Pos, g.Type.Ast(), g.Name}, nil
}

type ArrPattern struct {
	Pos      lexer.Position `"["`
	Element  *Pattern       `( @@`
	Elements []*Pattern     `  ( "," @@ )*`
	Rest     *Pattern       `  ( ";" @@ )? )? "]"`
}

func (g *ArrPattern) Ast() (expressions.Pattern, error) {
	var elPatterns []expressions.Pattern
	if g.Element != nil {
		elPatterns = make([]expressions.Pattern, len(g.Elements)+1)
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
	var restPattern expressions.Pattern
	if g.Rest == nil {
		restPattern = expressions.TypePattern{
			Pos:  g.Pos,
			Type: &types.Arr{types.Void{}},
		}
	} else {
		var err error
		restPattern, err = g.Rest.Ast()
		if err != nil {
			return nil, err
		}
	}
	return &expressions.ArrPattern{
		Pos:             g.Pos,
		ElementPatterns: elPatterns,
		RestPattern:     restPattern,
	}, nil
}

type ObjPattern struct {
	Pos    lexer.Position `"{"`
	Prop   *string        `( ( @Lid | @Op1 | @Op2 | @NumLiteral )`
	Value  *Pattern       `  ":" @@`
	Props  []string       `   ( "," ( @Lid | @Op1 | @Op2 | @NumLiteral )`
	Values []*Pattern     `     ":" @@ )* )? "}"`
}

func (g *ObjPattern) Ast() (expressions.Pattern, error) {
	propPatternMap := make(map[string]expressions.Pattern)
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
	return &expressions.ObjPattern{
		Pos:            g.Pos,
		PropPatternMap: propPatternMap,
	}, nil
}
