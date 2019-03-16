package parser

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/grammar"
)

func Parse(input string) (expressions.Expression, error) {
	parser, err := participle.Build(
		&grammar.Composition{},
		participle.Lexer(grammar.LexerDefinition),
		participle.Unquote("Str"),
		participle.Map(grammar.ToKeyword, "Prop"),
		participle.UseLookahead(0),
	)
	if err != nil {
		if lexerError, ok := err.(*lexer.Error); ok {
			return nil, errors.E(
				errors.Kind(errors.Syntax),
				errors.Pos(lexerError.Pos),
				errors.Message(lexerError.Message),
			)
		}
		return nil, err
	}
	composition := &grammar.Composition{}
	err = parser.ParseString(input, composition)
	if err != nil {
		if lexerError, ok := err.(*lexer.Error); ok {
			return nil, errors.E(
				errors.Kind(errors.Syntax),
				errors.Pos(lexerError.Pos),
				errors.Message(lexerError.Message),
			)
		}
		return nil, err
	}
	return composition.Ast(), nil
}
