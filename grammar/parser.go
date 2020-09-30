package grammar

import (
	"github.com/alecthomas/participle"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/types"
)

var parser *participle.Parser

var typeParser *participle.Parser

func init() {
	var err error
	parser, err = participle.Build(
		&Composition{},
		participle.Lexer(LexerDefinition),
		participle.Unquote("Str"),
	)
	if err != nil {
		panic(err)
	}
	typeParser, err = participle.Build(
		&Type{},
		participle.Lexer(LexerDefinition),
		participle.Unquote("Str"),
	)
	if err != nil {
		panic(err)
	}
}

func ParseComposition(input string) (expressions.Expression, error) {
	composition := &Composition{}
	err := parser.ParseString(input, composition)
	if err != nil {
		if parserError, ok := err.(participle.Error); ok {
			return nil, errors.E(
				errors.Code(errors.Syntax),
				errors.Pos(parserError.Token().Pos),
				errors.Message(parserError.Message()))

		}
		return nil, err
	}
	return composition.Ast()
}

func ParseType(input string) (types.Type, error) {
	t := &Type{}
	err := typeParser.ParseString(input, t)
	if err != nil {
		if parserError, ok := err.(participle.Error); ok {
			return nil, errors.E(
				errors.Code(errors.Syntax),
				errors.Pos(parserError.Token().Pos),
				errors.Message(parserError.Message()))

		}
		return nil, err
	}
	return t.Ast(), nil
}
