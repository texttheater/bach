package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

type Conditional struct {
	Pos          lexer.Position
	Pattern      *Pattern       `( "is" @@`
	Guard        *Composition   `  ( "with" @@)?`
	Condition    *Composition   `| "if" @@ )`
	Consequent   *Composition   `"then" @@` // long form
	Alternatives []*Alternative `( @@ )*`
	Alternative  *Composition   `( "else" @@ )? "ok"`
}

type Alternative struct {
	Pos        lexer.Position
	Pattern    *Pattern     `( "elseIs" @@`
	Guard      *Composition `  ( "with" @@ )?`
	Condition  *Composition `| "elseIf" @@ )`
	Consequent *Composition `"then" @@`
}

func (g *Conditional) Ast() (functions.Expression, error) {
	var pattern functions.Pattern
	var guard functions.Expression
	var err error
	if g.Pattern == nil {
		pattern = functions.TypePattern{
			Pos:  g.Pos,
			Type: types.AnyType{},
		}
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
	alternativePatterns := make([]functions.Pattern, len(g.Alternatives))
	alternativeGuards := make([]functions.Expression, len(g.Alternatives))
	alternativeConsequents := make([]functions.Expression, len(g.Alternatives))
	for i, alternative := range g.Alternatives {
		if alternative.Pattern == nil {
			alternativePatterns[i] = functions.TypePattern{
				Pos:  alternative.Pos,
				Type: types.AnyType{},
				Name: nil,
			}
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
	var alternative functions.Expression
	if g.Alternative != nil {
		alternative, err = g.Alternative.Ast()
		if err != nil {
			return nil, err
		}
	}
	return &functions.ConditionalExpression{
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
