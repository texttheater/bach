package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
)

type Conditional struct {
	Pos             lexer.Position
	Condition       *Composition   `"if" @@`
	Consequent      *Composition   `"then" @@`
	ElifConditions  []*Composition `( "elif" @@`
	ElifConsequents []*Composition `  "then" @@ )*`
	Alternative     *Composition   `"else" @@ "ok"`
}

func (g *Conditional) Ast() (expressions.Expression, error) {
	elifConditions := make([]expressions.Expression, len(g.ElifConditions))
	elifConsequents := make([]expressions.Expression, len(g.ElifConsequents))
	for i := range g.ElifConditions {
		var err error
		elifConditions[i], err = g.ElifConditions[i].Ast()
		if err != nil {
			return nil, err
		}
		elifConsequents[i], err = g.ElifConsequents[i].Ast()
		if err != nil {
			return nil, err
		}
	}
	condition, err := g.Condition.Ast()
	if err != nil {
		return nil, err
	}
	consequent, err := g.Consequent.Ast()
	if err != nil {
		return nil, err
	}
	alternative, err := g.Alternative.Ast()
	if err != nil {
		return nil, err
	}
	return &expressions.ConditionalExpression{
		Pos:             g.Pos,
		Condition:       condition,
		Consequent:      consequent,
		ElifConditions:  elifConditions,
		ElifConsequents: elifConsequents,
		Alternative:     alternative,
	}, nil
}
