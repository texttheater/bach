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
	interpreter.TestProgramStr(
		`for Arr<Num...> def f Num as =[a, b] a ok`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.NonExhaustiveMatch),
		),
		t,
	)
	interpreter.TestProgramStr(
		`for Obj<a: Num, Num> def f Num as ={a: a, b: b} a +b ok`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.NonExhaustiveMatch),
		),
		t,
	)
	interpreter.TestProgramStr(
		`for Obj<a: Num, Num> def f Num as ={b: b} b ok`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.NonExhaustiveMatch),
		),
		t,
	)
}
