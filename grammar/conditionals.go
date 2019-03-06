package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/ast"
)

type Conditional struct {
	Pos             lexer.Position
	Condition       *Composition   `"if" @@`
	Consequent      *Composition   `"then" @@`
	ElifConditions  []*Composition `{ "elif" @@`
	ElifConsequents []*Composition `  "then" @@ }`
	Alternative     *Composition   `"else" @@ "ok"`
}

func (g *Conditional) Ast() ast.Expression {
	elifConditions := make([]ast.Expression, len(g.ElifConditions))
	elifConsequents := make([]ast.Expression, len(g.ElifConsequents))
	for i := range g.ElifConditions {
		elifConditions[i] = g.ElifConditions[i].Ast()
		elifConsequents[i] = g.ElifConsequents[i].Ast()
	}
	return &ast.ConditionalExpression{
		Pos:             g.Pos,
		Condition:       g.Condition.Ast(),
		Consequent:      g.Consequent.Ast(),
		ElifConditions:  elifConditions,
		ElifConsequents: elifConsequents,
		Alternative:     g.Alternative.Ast(),
	}
}
