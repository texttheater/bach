package grammar

import (
	"regexp"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
)

type Regexp struct {
	Pos    lexer.Position
	Regexp string `@Regexp`
}

func (g *Regexp) Ast() (functions.Expression, error) {
	regexpString := g.Regexp[1 : len(g.Regexp)-1]
	regexp, err := regexp.Compile(regexpString)
	if err != nil {
		return nil, errors.E(
			errors.Pos(g.Pos),
			errors.Code(errors.BadRegexp),
			errors.Message(err.Error()),
		)
	}
	regexpExpression := &functions.RegexpExpression{
		Pos:    g.Pos,
		Regexp: regexp,
	}
	return regexpExpression, nil
}
