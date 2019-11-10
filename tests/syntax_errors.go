package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
)

func TestSyntaxErrors(t *testing.T) {
	TestProgram(
		"&",
		nil,
		nil,
		errors.E(
			errors.Code(errors.Syntax),
		),
		t,
	)
	TestProgram(
		"drop",
		nil,
		nil,
		errors.E(
			errors.Code(errors.Syntax),
		),
		t,
	)
	TestProgram(
		"if true then drop else true ok",
		nil,
		nil,
		errors.E(
			errors.Code(errors.Syntax),
		),
		t,
	)
}
