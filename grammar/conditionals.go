package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/types"
)

type Conditional struct {
	Pos              lexer.Position
	Pattern          *Pattern           `( "is" @@`
	Guard            *Composition       `  ( "with" @@)?`
	Condition        *Composition       `| "if" @@ )`
	Consequent       *Composition       `( "then" @@`
	Alternatives     []*Alternative     `  ( @@ )*`
	Alternative      *Composition       `  ( "else" @@ )? "ok"`
	PredAlternatives []*PredAlternative `| @@+ )?`
}

type Alternative struct {
	Pos        lexer.Position
	Pattern    *Pattern     `( "elis" @@`
	Guard      *Composition `  ( "with" @@ )?`
	Condition  *Composition `| "elif" @@ )`
	Consequent *Composition `"then" @@`
}

type PredAlternative struct {
	Pos       lexer.Position
	Pattern   *Pattern     `( "elis" @@`
	Guard     *Composition `  ( "with" @@ )?`
	Condition *Composition `| "elif" @@ )`
}

func (g *Conditional) Ast() (expressions.Expression, error) {
	var pattern expressions.Pattern
	var guard expressions.Expression
	var err error
	// build pattern and guard
	if g.Pattern == nil {
		pattern = expressions.TypePattern{
			Pos:  g.Pos,
			Type: types.Any{},
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
	// build consequent and alternatives
	var consequent expressions.Expression
	var alternativePatterns []expressions.Pattern
	var alternativeGuards []expressions.Expression
	var alternativeConsequents []expressions.Expression
	if g.Consequent == nil { // predicate form
		alternativePatterns = make([]expressions.Pattern, len(g.PredAlternatives))
		alternativeGuards = make([]expressions.Expression, len(g.PredAlternatives))
		alternativeConsequents = make([]expressions.Expression, len(g.PredAlternatives))
		yesExpression := &expressions.WrapExpression{
			Pos:  g.Pos,
			Prop: "yes",
		}
		consequent = yesExpression
		// build alternative patterns, guards, and consequents
		for i, alternative := range g.PredAlternatives {
			if alternative.Pattern == nil {
				alternativePatterns[i] = expressions.TypePattern{
					Pos:  alternative.Pos,
					Type: types.Any{},
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
			alternativeConsequents[i] = yesExpression
		}
	} else {
		alternativePatterns = make([]expressions.Pattern, len(g.Alternatives))
		alternativeGuards = make([]expressions.Expression, len(g.Alternatives))
		alternativeConsequents = make([]expressions.Expression, len(g.Alternatives))
		consequent, err = g.Consequent.Ast()
		if err != nil {
			return nil, err
		}
		// build alternative patterns, guards, and consequents
		for i, alternative := range g.Alternatives {
			if alternative.Pattern == nil {
				alternativePatterns[i] = expressions.TypePattern{
					Pos:  alternative.Pos,
					Type: types.Any{},
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
	}
	// build final alternative
	var alternative expressions.Expression
	if g.Consequent == nil { // predicate form
		alternative = &expressions.WrapExpression{
			Pos:  g.Pos,
			Prop: "no",
		}
	} else {
		if g.Alternative != nil {
			alternative, err = g.Alternative.Ast()
		}
		if err != nil {
			return nil, err
		}
	}
	// return
	return &expressions.ConditionalExpression{
		Pos:                    g.Pos,
		Pattern:                pattern,
		Guard:                  guard,
		Consequent:             consequent,
		AlternativePatterns:    alternativePatterns,
		AlternativeGuards:      alternativeGuards,
		AlternativeConsequents: alternativeConsequents,
		Alternative:            alternative,
	}, nil
}
