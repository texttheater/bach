package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/patterns"
	"github.com/texttheater/bach/types"
)

type Match struct {
	Pos          lexer.Position
	Pattern      *Pattern       `( "is" @@`
	Guard        *Composition   `  ( "with" @@)?`
	Condition    *Composition   `| "if" @@ )`
	Consequent   *Composition   `"then" @@`
	Alternatives []*Alternative `( @@ )*`
	Alternative  *Composition   `( "else" @@ )? "ok"`
}

type Alternative struct {
	Pos        lexer.Position
	Pattern    *Pattern     `( "elis" @@`
	Guard      *Composition `  ( "with" @@ )?`
	Condition  *Composition `| "elif" @@ )`
	Consequent *Composition `"then" @@`
}

func (g *Match) Ast() (expressions.Expression, error) {
	var pattern patterns.Pattern
	var guard expressions.Expression
	var err error
	if g.Pattern == nil {
		pattern = patterns.TypePattern{g.Pos, types.AnyType{}, nil}
		guard, err = g.Condition.Ast()
		if err != nil {
			return nil, err
		}
	} else {
		pattern, err = g.Pattern.Ast()
		if err != nil {
			return nil, err
		}
		if g.Guard != nil {
			guard, err = g.Guard.Ast()
			if err != nil {
				return nil, err
			}
		}
	}
	consequent, err := g.Consequent.Ast()
	if err != nil {
		return nil, err
	}
	alternativePatterns := make([]patterns.Pattern, len(g.Alternatives))
	alternativeGuards := make([]expressions.Expression, len(g.Alternatives))
	alternativeConsequents := make([]expressions.Expression, len(g.Alternatives))
	for i, alternative := range g.Alternatives {
		if alternative.Pattern == nil {
			alternativePatterns[i] = patterns.TypePattern{alternative.Pos, types.AnyType{}, nil}
			alternativeGuards[i], err = alternative.Condition.Ast()
			if err != nil {
				return nil, err
			}
		} else {
			alternativePatterns[i], err = alternative.Pattern.Ast()
			if err != nil {
				return nil, err
			}
			if alternative.Guard != nil {
				alternativeGuards[i], err = alternative.Guard.Ast()
				if err != nil {
					return nil, err
				}
			}
		}
		alternativeConsequents[i], err = alternative.Consequent.Ast()
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
		ElisPatterns:    alternativePatterns,
		ElisGuards:      alternativeGuards,
		ElisConsequents: alternativeConsequents,
		Alternative:     alternative,
	}, nil
}
