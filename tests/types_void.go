package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
)

func TestVoidType(t *testing.T) {
	TestProgram(
		`[] each drop 0 all`,
		nil,
		nil,
		states.E(
			states.Code(states.ComposeWithVoid)),

		t,
	)
}
