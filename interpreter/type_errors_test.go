package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/types"
)

func TestTypeErrors(t *testing.T) {
	interpreter.TestProgramStr(
		`3 <2 +1`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.NoSuchFuncer),
			errors.InputType(types.Bool{}),
			errors.Name("+"),
			errors.NumParams(1),
		),
		t,
	)
	interpreter.TestProgramStr(
		`+`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.NoSuchFuncer),
			errors.InputType(types.Null{}),
			errors.Name("+"),
			errors.NumParams(0),
		),
		t,
	)
	interpreter.TestProgramStr(
		`hurz`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.NoSuchFuncer),
		),
		t,
	)
}
