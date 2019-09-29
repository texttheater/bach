package grammar

import (
	"regexp"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
)

type RegexpFindFirst struct {
	Pos             lexer.Position
	RegexpFindFirst *string `@RegexpFindFirst`
}

func (g *RegexpFindFirst) Ast() (functions.Expression, error) {
	regexpString := (*g.RegexpFindFirst)[2 : len(*g.RegexpFindFirst)-1]
	regexp, err := regexp.Compile(regexpString)
	if err != nil {
		return nil, err
	}
	regexpFindFirstExpression := &functions.RegexpFindFirstExpression{
		Pos:    g.Pos,
		Regexp: regexp,
	}
	return regexpFindFirstExpression, nil
}

type RegexpFindAll struct {
	Pos           lexer.Position
	RegexpFindAll *string `@RegexpFindAll`
}

func (g *RegexpFindAll) Ast() (functions.Expression, error) {
	regexpString := (*g.RegexpFindAll)[2 : len(*g.RegexpFindAll)-2]
	regexp, err := regexp.Compile(regexpString)
	if err != nil {
		return nil, err
	}
	regexpFindAllExpression := &functions.RegexpFindAllExpression{
		Pos:    g.Pos,
		Regexp: regexp,
	}
	return regexpFindAllExpression, nil
}
