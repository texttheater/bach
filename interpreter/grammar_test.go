package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestGrammar(t *testing.T) {
	interpreter.TestProgram(
		"&",
		nil,
		nil,
		errors.SyntaxError(
			errors.Code(errors.Syntax),
		),
		t,
	)
	// The following program requires a lookahead of 1 (participle's
	// default) so > is not interpreted as a property identifier.
	interpreter.TestProgram(
		"for Str def f Obj<> as {} ok",
		types.Null{},
		states.NullValue{},
		nil,
		t,
	)
}
