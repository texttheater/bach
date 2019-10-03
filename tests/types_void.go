package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
)

func TestVoidType(t *testing.T) {
	TestProgram(
		`drop 0`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.ComposeWithVoid),
		),
		t,
	)
	TestProgram(
		`drop`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.VoidProgram),
		),
		t,
	)
}
