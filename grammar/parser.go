package grammar

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
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
		if lexerError, ok := err.(*lexer.Error); ok {
			return nil, errors.E(
				errors.Code(errors.Syntax),
				errors.Pos(lexerError.Pos),
				errors.Message(lexerError.Message),
			)
		}
		return nil, err
	}
	composition := &Composition{}
	err = parser.ParseString(input, composition)
	if err != nil {
		if lexerError, ok := err.(*lexer.Error); ok {
			return nil, errors.E(
				errors.Code(errors.Syntax),
				errors.Pos(lexerError.Pos),
				errors.Message(lexerError.Message),
			)
		}
		return nil, err
	}
	return composition.Ast()
}
