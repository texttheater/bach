package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/types"
)

type Conditional struct {
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

func (g *Conditional) Ast() (expressions.Expression, error) {
	var pattern expressions.Pattern
	var guard expressions.Expression
	var err error
	if g.Pattern == nil {
		pattern = expressions.TypePattern{
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
	var consequent expressions.Expression
	consequent, err = g.Consequent.Ast()
	if err != nil {
		return nil, err
	}
	alternativePatterns := make([]expressions.Pattern, len(g.Alternatives))
	alternativeGuards := make([]expressions.Expression, len(g.Alternatives))
	alternativeConsequents := make([]expressions.Expression, len(g.Alternatives))
	for i, alternative := range g.Alternatives {
		if alternative.Pattern == nil {
			alternativePatterns[i] = expressions.TypePattern{
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
	var alternative expressions.Expression
	if g.Alternative != nil {
		alternative, err = g.Alternative.Ast()
	}
	if err != nil {
		return nil, err
	}
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
