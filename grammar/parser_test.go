package grammar_test

import (
	"testing"

	"github.com/alecthomas/participle"
	"github.com/texttheater/bach/grammar"
)

func TestNumGetter(t *testing.T) {
	parser, err := participle.Build(
		&grammar.Getter{},
		participle.Lexer(grammar.LexerDefinition),
	)
	if err != nil {
		t.Fatal(err)
	}
	getter := &grammar.Getter{}
	err = parser.ParseString("@-1", getter)
	if err != nil {
		t.Fatal(err)
	}
	if getter.NumGetter == nil {
		t.Fatal("@-1 not parsed as NumGetter")
	}
}

func TestFilter(t *testing.T) {
	parser, err := participle.Build(
		&grammar.Filter{},
		participle.Lexer(grammar.LexerDefinition),
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
