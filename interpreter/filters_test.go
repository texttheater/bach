package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
)

func TestFilters(t *testing.T) {
	interpreter.TestProgramStr(
		`[1, 2, 3] each(if ==1 then "a" elif ==2 then "b" else "c" ok)`,
		`Arr<Str>`,
		`["a", "b", "c"]`,
		nil,
		t,
	)
	interpreter.TestProgram(
		`[1, 2, 3] each(if ==1 then "a" elif ==2 then "b" else "c")`,
		nil,
		nil,
		errors.SyntaxError(
			errors.Code(errors.Syntax),
		),
		t,
	)
	interpreter.TestProgramStr(
		`[1, 2, 3] each(+1)`,
		`Arr<Num>`,
		`[2, 3, 4]`,
		nil,
		t,
	)
}
