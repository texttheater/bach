package grammar

import (
	"github.com/alecthomas/participle"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
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
			return nil, states.E(
				states.Code(states.Syntax),
				states.Pos(parserError.Token().Pos),
				states.Message(parserError.Message()))

		}
		return nil, err
	}
	composition := &Composition{}
	err = parser.ParseString(input, composition)
	if err != nil {
		if parserError, ok := err.(participle.Error); ok {
			return nil, states.E(
				states.Code(states.Syntax),
				states.Pos(parserError.Token().Pos),
				states.Message(parserError.Message()))

		}
		return nil, err
	}
	return composition.Ast()
}
