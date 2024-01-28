package tests_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/tests"
	"github.com/texttheater/bach/types"
)

func TestTypeErrors(t *testing.T) {
	tests.TestProgramStr(
		`3 <2 +1`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.NoSuchFunction),
			errors.InputType(types.Bool{}),
			errors.Name("+"),
			errors.NumParams(1),
		),
		t,
	)
	tests.TestProgramStr(
		`+`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.NoSuchFunction),
			errors.InputType(types.Null{}),
			errors.Name("+"),
			errors.NumParams(0),
		),
		t,
	)
	tests.TestProgramStr(
		`hurz`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.NoSuchFunction),
		),
		t,
	)
}
