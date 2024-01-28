package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
)

func TestFilters(t *testing.T) {
	interpreter.TestProgramStr(
		`["a", 1, "b", 2, "c", 3] keep(is Num with %2 >0 elis Str)`,
		`Arr<Num|Str>`,
		`["a", 1, "b", "c", 3]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[{n: 1}, {n: 2}, {n: 3}] keep(is {n: n} with n %2 >0)`,
		`Arr<Obj<n: Num, Void>>`,
		`[{n: 1}, {n: 3}]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[1, 2, 3, 4, 5, 6] keep(if %2 ==0) each(*2)`,
		`Arr<Num>`,
		`[4, 8, 12]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[1, 2, 3, 4, 5, 6] keep(if %2 ==0 not) each(id)`,
		`Arr<Num>`,
		`[1, 3, 5]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[1, 2, 3] keep(if %2 ==0 not) each(+1)`,
		`Arr<Num>`,
		`[2, 4]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[1, 2, 3] keep(if false)`,
		`Arr<Num>`,
		`[]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[{n: 1}, 2, {n: 3}] keep(is {n: n}) each(@n)`,
		`Arr<Num>`,
		`[1, 3]`,
		nil,
		t,
	)
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
