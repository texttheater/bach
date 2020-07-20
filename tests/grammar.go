package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestGrammar(t *testing.T) {
	TestProgram(
		"&",
		nil,
		nil,
		states.E(
			states.Code(states.Syntax)),

		t,
	)
	TestProgram(
		"drop",
		nil,
		nil,
		states.E(
			states.Code(states.Syntax)),

		t,
	)
	TestProgram(
		"if true then drop else true ok",
		nil,
		nil,
		states.E(
			states.Code(states.Syntax)),

		t,
	)
	// The following program requires a lookahead of 1 (participle's
	// default) so > is not interpreted as a property identifier.
	TestProgram(
		"for Str def f Obj<> as {} ok",
		types.NullType{},
		states.NullValue{},
		nil,
		t,
	)
}
