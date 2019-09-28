package grammar

import (
	"regexp/syntax"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
)

type RegexpMatch struct {
	Pos         lexer.Position
	RegexpMatch *string `@RegexpMatch`
}

func (g *RegexpMatch) Ast() (functions.Expression, error) {
	regexpString := (*g.RegexpMatch)[2 : len(*g.RegexpMatch)-1]
	regexp, err := syntax.Parse(regexpString, 0)
	if err != nil {
		return nil, err
	}
	regexpMatchExpression := &functions.RegexpMatchExpression{
		Pos:    g.Pos,
		Regexp: regexp,
	}
	return regexpMatchExpression, nil
}
