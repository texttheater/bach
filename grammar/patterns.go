package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/patterns"
	"github.com/texttheater/bach/types"
)

type Pattern struct {
	Pos        lexer.Position
	Type       *Type       `  @@`
	ArrPattern *ArrPattern `| @@`
}

func (g *Pattern) Ast() (patterns.Pattern, error) {
	var t types.Type
	if g.Type != nil {
		t = g.Type.Ast()
	} else if g.ArrPattern != nil {
		p, err := g.ArrPattern.Ast()
		if err != nil {
			return nil, err
		}
		return p, nil
	} else {
		panic("invalid pattern")
	}
	return patterns.TypePattern{g.Pos, t}, nil
}

type ArrPattern struct {
	Pos      lexer.Position `"["`
	Element  *Pattern       `( @@`
	Elements []*Pattern     `  ( "," @@ )* )? "]"`
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
	return &patterns.ArrPattern{g.Pos, elPatterns}, nil
}
