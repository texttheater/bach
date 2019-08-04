package grammar

import (
	"regexp/syntax"
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
)

type RegexpCall struct {
	Pos        lexer.Position
	NameRegexp *string `@NameRegexp`
}

func (g *RegexpCall) Ast() (expressions.Expression, error) {
	index := strings.Index(*g.NameRegexp, "/")
	name := (*g.NameRegexp)[:index]
	regexpString := (*g.NameRegexp)[index : len(*g.NameRegexp)-1]
	regexp, err := syntax.Parse(regexpString, 0)
	if err != nil {
		return nil, err
	}
	regexpExpression := &expressions.RegexpExpression{
		Pos:    g.Pos,
		Regexp: regexp,
	}
	callExpression := &expressions.CallExpression{
		Pos:  g.Pos,
		Name: name,
		Args: []expressions.Expression{
			regexpExpression,
		},
	}
	return callExpression, nil
}
