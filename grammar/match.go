package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/patterns"
)

type Match struct {
	Pos             lexer.Position
	Pattern         *Pattern       `"is" @@`
	Guard           *Composition   `( "with" @@)?`
	Consequent      *Composition   `"then" @@`
	ElisPatterns    []*Pattern     `( "elis" @@`
	ElisGuards      []*Composition `  ( "with" @@)?`
	ElisConsequents []*Composition `  "then" @@ )*`
	Alternative     *Composition   `( "else" @@ )? "ok"`
}

func (g *Match) Ast() (expressions.Expression, error) {
	pattern, err := g.Pattern.Ast()
	if err != nil {
		return nil, err
	}
	consequent, err := g.Consequent.Ast()
	if err != nil {
		return nil, err
	}
	elisPatterns := make([]patterns.Pattern, len(g.ElisPatterns))
	elisConsequents := make([]expressions.Expression, len(g.ElisConsequents))
	for i := range g.ElisPatterns {
		elisPatterns[i], err = g.ElisPatterns[i].Ast()
		if err != nil {
			return nil, err
		}
		elisConsequents[i], err = g.ElisConsequents[i].Ast()
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
		Consequent:      consequent,
		ElisPatterns:    elisPatterns,
		ElisConsequents: elisConsequents,
		Alternative:     alternative,
	}, nil
}
