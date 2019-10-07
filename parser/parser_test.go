package parser_test

import (
	"github.com/alecthomas/participle"
	"github.com/texttheater/bach/grammar"

	"testing"
)

func TestFilter(t *testing.T) {
	parser, err := participle.Build(
		&grammar.Filter{},
		participle.Lexer(grammar.LexerDefinition),
		participle.Unquote("Str"),
		participle.Map(grammar.ToKeyword, "Lid"),
		participle.UseLookahead(0),
	)
	if err != nil {
		t.Fatal(err)
	}
	filter := &grammar.Filter{}
	err = parser.ParseString("each if %2 >0 elis Str all", filter)
	if err != nil {
		t.Fatal(err)
	}
}
