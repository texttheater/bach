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

func TestTypes(t *testing.T) {
	_, err := grammar.ParseType("Obj<>")
	if err != nil {
		t.Fatal(err)
	}
	_, err = grammar.ParseType("Obj<a: Bool>")
	if err != nil {
		t.Fatal(err)
	}
	_, err = grammar.ParseType("Obj<a: Bool, b: Bool, Bool>")
	if err != nil {
		t.Fatal(err)
	}
	_, err = grammar.ParseType("Obj<a: Num, b: Num>")
	if err != nil {
		t.Fatal(err)
	}
}

func TestTypeTemplates(t *testing.T) {
	_, err := grammar.ParseTypeTemplate("Obj<>")
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	_, err = grammar.ParseTypeTemplate("Obj<a: Bool, b: Bool, Bool>")
	if err != nil {
		t.Fatal(err)
	}
	_, err = grammar.ParseTypeTemplate("Obj<a: Num, b: Num>")
	if err != nil {
		t.Fatal(err)
	}
}
