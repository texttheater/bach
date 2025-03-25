package docutil_test

import (
	"encoding/json"
	"testing"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/docutil"
	"github.com/texttheater/bach/errors"
)

func TestNullError(t *testing.T) {
	input := "null"
	e, err := docutil.ParseError(input)
	if err != nil {
		t.Errorf(`failed to parse %s: %s`, input, err)
		return
	}
	if e != nil {
		t.Errorf(`failed to parse %s: expected %s, got %s`, input, error(nil), e)
	}
}

func TestEmptyError(t *testing.T) {
	input := "{}"
	e, err := docutil.ParseError(input)
	if err != nil {
		t.Errorf(`failed to parse %s: %s`, input, err)
		return
	}
	if e != nil {
		t.Errorf(`failed to parse %s: expected %s, got %s`, input, error(nil), e)
	}
}

func TestNonemptyError(t *testing.T) {
	input := `{"Code": "ArgHasWrongOutputType"}`
	var e errors.E
	err := json.Unmarshal([]byte(input), &e)
	if err != nil {
		t.Errorf(`failed to parse %s: %s`, input, err)
		return
	}
	if e.Pos != nil {
		t.Errorf(`failed to parse %s: expected Pos %s, got %s`, input, (*lexer.Position)(nil), e.Pos)
	}
}
