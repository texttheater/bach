package grammar

import (
	"github.com/alecthomas/participle"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
)

func Parse(input string) (functions.Expression, error) {
	parser, err := participle.Build(
		&Composition{},
		participle.Lexer(LexerDefinition),
		participle.Unquote("Str"),
		participle.Map(ToKeyword, "Lid"),
	)
	if err != nil {
		if parserError, ok := err.(participle.Error); ok {
			return nil, errors.E(
				errors.Code(errors.Syntax),
				errors.Pos(parserError.Token().Pos),
				errors.Message(parserError.Message()),
			)
		}
		return nil, err
	}
	composition := &Composition{}
	err = parser.ParseString(input, composition)
	if err != nil {
		if parserError, ok := err.(participle.Error); ok {
			return nil, errors.E(
				errors.Code(errors.Syntax),
				errors.Pos(parserError.Token().Pos),
				errors.Message(parserError.Message()),
			)
		}
		return nil, err
	}
	return composition.Ast()
}
