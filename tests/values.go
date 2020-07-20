package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestValues(t *testing.T) {
	TestProgram(
		`false if id then 1 else fatal ok`,
		types.NumType{},
		nil,
		errors.E(
			errors.Code(errors.UnexpectedValue),
			errors.GotValue(states.BoolValue(false))),

		t,
	)
}
