package grammar

import (
	"github.com/alecthomas/participle"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/types"
)

var parser *participle.Parser

var typeParser *participle.Parser

var templateParser *participle.Parser

var paramParser *participle.Parser

func init() {
	var err error
	parser, err = participle.Build(
		&Composition{},
		participle.Lexer(LexerDefinition),
	)
	if err != nil {
		panic(err)
	}
	typeParser, err = participle.Build(
		&Type{},
		participle.Lexer(LexerDefinition),
	)
	if err != nil {
		panic(err)
	}
	templateParser, err = participle.Build(
		&TypeTemplate{},
		participle.Lexer(LexerDefinition),
	)
	if err != nil {
		panic(err)
	}
	paramParser, err = participle.Build(
		&Param{},
		participle.Lexer(LexerDefinition),
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
			return nil, errors.SyntaxError(
				errors.Code(errors.Syntax),
				errors.Pos(parserError.Token().Pos),
				errors.Message(parserError.Message()),
			)
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
			return nil, errors.SyntaxError(
				errors.Code(errors.Syntax),
				errors.Pos(parserError.Token().Pos),
				errors.Message(parserError.Message()),
			)
		}
		return nil, err
	}
	return t.Ast(), nil
}

func ParseTypeTemplate(input string) (types.Type, error) {
	t := &TypeTemplate{}
	err := templateParser.ParseString(input, t)
	if err != nil {
		if parserError, ok := err.(participle.Error); ok {
			return nil, errors.SyntaxError(
				errors.Code(errors.Syntax),
				errors.Pos(parserError.Token().Pos),
				errors.Message(parserError.Message()),
			)
		}
		return nil, err
	}
	return t.Ast(), nil
}

func ParseParam(input string) (*params.Param, error) {
	p := &Param{}
	err := paramParser.ParseString(input, p)
	if err != nil {
		if parserError, ok := err.(participle.Error); ok {
			return nil, errors.SyntaxError(
				errors.Code(errors.Syntax),
				errors.Pos(parserError.Token().Pos),
				errors.Message(parserError.Message()),
			)
		}
		return nil, err
	}
	return p.Ast()
}
