package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/types"
)

type Assignment struct {
	Pos            lexer.Position
	NameAssignment *NameAssignment `  @@`
	ArrAssignment  *ArrAssignment  `| @@`
}

func (g *Assignment) Ast() (expressions.Expression, error) {
	if g.NameAssignment != nil {
		return g.NameAssignment.Ast()
	}
	if g.ArrAssignment != nil {
		return g.ArrAssignment.Ast()
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
