package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/patterns"
)

type Match struct {
	Pos         lexer.Position
	Pattern     *Pattern     `"is" @@`
	Guard       *Composition `( "with" @@)?`
	Consequent  *Composition `"then" @@`
	Eliss       []*Elis      `( @@ )*`
	Alternative *Composition `( "else" @@ )? "ok"`
}

type Elis struct {
	Pos        lexer.Position
	Pattern    *Pattern     `"elis" @@`
	Guard      *Composition `( "with" @@ )?`
	Consequent *Composition `"then" @@`
}

func (g *Match) Ast() (expressions.Expression, error) {
	pattern, err := g.Pattern.Ast()
	if err != nil {
		return nil, err
	}
	var guard expressions.Expression
	if g.Guard != nil {
		guard, err = g.Guard.Ast()
		if err != nil {
			return nil, err
		}
	}
	consequent, err := g.Consequent.Ast()
	if err != nil {
		return nil, err
	}
	elisPatterns := make([]patterns.Pattern, len(g.Eliss))
	elisGuards := make([]expressions.Expression, len(g.Eliss))
	elisConsequents := make([]expressions.Expression, len(g.Eliss))
	for i, elis := range g.Eliss {
		elisPatterns[i], err = elis.Pattern.Ast()
		if err != nil {
			return nil, err
		}
		if elis.Guard != nil {
			elisGuards[i], err = elis.Guard.Ast()
			if err != nil {
				return nil, err
			}
		}
		elisConsequents[i], err = elis.Consequent.Ast()
		if err != nil {
			return nil, err
		}
	}
	var alternative expressions.Expression
	if g.Alternative != nil {
		alternative, err = g.Alternative.Ast()
		if err != nil {
			return nil, err
		}
	}
	return &expressions.MatchExpression{
		Pos:             g.Pos,
		Pattern:         pattern,
		Guard:           guard,
		Consequent:      consequent,
		ElisPatterns:    elisPatterns,
		ElisGuards:      elisGuards,
		ElisConsequents: elisConsequents,
		Alternative:     alternative,
	}, nil
}
