package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/types"
)

// FIXME should not use patterns but another thing that doesn't accept types
// without names (does not make sense for assignments)

type Assignment struct {
	Pos            lexer.Position
	NameAssignment *NameAssignment `  @@`
	ArrAssignment  *ArrAssignment  `| @@`
	ObjAssignment  *ObjAssignment  `| @@`
}

func (g *Assignment) Ast() (expressions.Expression, error) {
	if g.NameAssignment != nil {
		return g.NameAssignment.Ast()
	}
	if g.ArrAssignment != nil {
		return g.ArrAssignment.Ast()
	}
	if g.ObjAssignment != nil {
		return g.ObjAssignment.Ast()
	}
	panic("invalid assignment")
}

type NameAssignment struct {
	Pos    lexer.Position
	EqName string `@EqName`
}

func (g *NameAssignment) Ast() (expressions.Expression, error) {
	name := g.EqName[1:]
	pattern := expressions.TypePattern{g.Pos, types.Any{}, &name}
	return &expressions.AssignmentExpression{
		Pos:     g.Pos,
		Pattern: pattern,
	}, nil
}

type ArrAssignment struct {
	Pos      lexer.Position `"=["`
	Element  *Pattern       `( @@`
	Elements []*Pattern     `  ( "," @@ )*`
	Rest     *Pattern       `  ( ";" @@ )? )? "]"`
}

func (g *ArrAssignment) Ast() (expressions.Expression, error) {
	numEls := 0
	if g.Element != nil {
		numEls = 1
		numEls += len(g.Elements)
	}
	elPatterns := make([]expressions.Pattern, numEls)
	if g.Element != nil {
		pattern, err := g.Element.Ast()
		if err != nil {
			return nil, err
		}
		elPatterns[0] = pattern
	}
	for i, element := range g.Elements {
		pattern, err := element.Ast()
		if err != nil {
			return nil, err
		}
		elPatterns[i+1] = pattern
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
	return &expressions.AssignmentExpression{
		Pos: g.Pos,
		Pattern: expressions.ArrPattern{
			Pos:             g.Pos,
			ElementPatterns: elPatterns,
			RestPattern:     restPattern,
		},
	}, nil
}

type ObjAssignment struct {
	Pos    lexer.Position `"={"`
	Prop   *string        `( ( @Lid | @Op1 | @Op2 | @NumLiteral )`
	Value  *Pattern       `  ":" @@`
	Props  []string       `   ( "," ( @Lid | @Op1 | @Op2 | @NumLiteral )`
	Values []*Pattern     `     ":" @@ )* )? "}"`
}

func (g *ObjAssignment) Ast() (expressions.Expression, error) {
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
	return &expressions.AssignmentExpression{
		Pos: g.Pos,
		Pattern: expressions.ObjPattern{
			Pos:            g.Pos,
			PropPatternMap: propPatternMap,
		},
	}, nil
}
