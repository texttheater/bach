package grammar

import (
	"regexp/syntax"
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
)

type RegexpCall struct {
	Pos        lexer.Position
	NameRegexp *string `@NameRegexp`
}

func (g *RegexpCall) Ast() (functions.Expression, error) {
	index := strings.Index(*g.NameRegexp, "/")
	name := (*g.NameRegexp)[:index]
	regexpString := (*g.NameRegexp)[index : len(*g.NameRegexp)-1]
	regexp, err := syntax.Parse(regexpString, 0)
	if err != nil {
		return nil, err
	}
	regexpExpression := &functions.RegexpExpression{
		Pos:    g.Pos,
		Regexp: regexp,
	}
	callExpression := &functions.CallExpression{
		Pos:  g.Pos,
		Name: name,
		Args: []functions.Expression{
			regexpExpression,
		},
	}
	return callExpression, nil
}
