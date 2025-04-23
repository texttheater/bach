package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
)

func TestAssignment(t *testing.T) {
	interpreter.TestProgramStr(
		`2 =Num`,
		``,
		``,
		errors.SyntaxError(
			errors.Code(errors.Syntax),
		),
		t,
	)
}
