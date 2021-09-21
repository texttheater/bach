package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/types"
)

func TestTypeErrors(t *testing.T) {
	TestProgramStr(
		`3 <2 +1`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.NoSuchFunction),
			errors.InputType(types.BoolType{}),
			errors.Name("+"),
			errors.NumParams(1),
		),
		t,
	)
	TestProgramStr(
		`+`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.NoSuchFunction),
			errors.InputType(types.NullType{}),
			errors.Name("+"),
			errors.NumParams(0),
		),
		t,
	)
	TestProgramStr(
		`hurz`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.NoSuchFunction),
		),
		t,
	)
}
