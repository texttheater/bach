package grammar

import (
	"regexp"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/expressions"
)

type Regexp struct {
	Pos    lexer.Position
	Regexp string `@Regexp`
}

func (g *Regexp) Ast() (expressions.Expression, error) {
	regexpString := g.Regexp[1 : len(g.Regexp)-1]
	regexp, err := regexp.Compile(regexpString)
	if err != nil {
		return nil, errors.SyntaxError(
			errors.Code(errors.BadRegexp),
			errors.Pos(g.Pos),
			errors.Message(err.Error()),
		)
	}
	regexpExpression := &expressions.RegexpExpression{
		Pos:    g.Pos,
		Regexp: regexp,
	}
	return regexpExpression, nil
}
