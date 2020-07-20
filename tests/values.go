package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestValues(t *testing.T) {
	TestProgram(
		`false if id then 1 else fatal ok`,
		types.NumType{},
		nil,
		states.E(
			states.Code(states.UnexpectedValue),
			states.GotValue(states.BoolValue(false))),

		t,
	)
}
