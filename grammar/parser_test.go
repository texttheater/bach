package grammar_test

import (
	"testing"

	"github.com/alecthomas/participle"
	"github.com/texttheater/bach/grammar"
)

func TestFilter(t *testing.T) {
	parser, err := participle.Build(
		&grammar.Filter{},
		participle.Lexer(grammar.LexerDefinition),
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

func TestTypes(t *testing.T) {
	_, err := grammar.ParseType("Obj<>")
	if err != nil {
		t.Fatal(err)
	}
}
