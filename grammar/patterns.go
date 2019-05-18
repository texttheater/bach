package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/patterns"
	"github.com/texttheater/bach/types"
)

type Pattern struct {
	Pos        lexer.Position
	NullType   *NullType   `  @@`
	ReaderType *ReaderType `| @@`
	BoolType   *BoolType   `| @@`
	NumType    *NumType    `| @@`
	StrType    *StrType    `| @@`
	SeqType    *SeqType    `| @@`
	ArrPattern *ArrPattern `| @@`
}

func (g *Pattern) Ast() (patterns.Pattern, error) {
	var t types.Type
	if g.NullType != nil {
		t = g.NullType.Ast()
	} else if g.ReaderType != nil {
		t = g.ReaderType.Ast()
	} else if g.BoolType != nil {
		t = g.BoolType.Ast()
	} else if g.NumType != nil {
		t = g.NumType.Ast()
	} else if g.StrType != nil {
		t = g.StrType.Ast()
	} else if g.SeqType != nil {
		var err error
		t, err = g.SeqType.Ast()
		if err != nil {
			return nil, err
		}
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
