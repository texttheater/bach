package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
)

func TestVoidType(t *testing.T) {
	TestProgram(
		`[] each drop 0 all`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.ComposeWithVoid)),

		t,
	)
}
