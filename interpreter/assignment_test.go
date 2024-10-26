package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
)

func TestAssignment(t *testing.T) {
	interpreter.TestProgramStr(
		`1 +1 =a 3 *2 +a`,
		`Num`,
		`8`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`1 +1 ==2 =p 1 +1 ==1 =q p ==q not`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[1, 2, 3] =[a, b, c] a`,
		`Num`,
		`1`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[1, 2, 3] =[a, b, c] c`,
		`Num`,
		`3`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[1, 2, 3] =[a;r] r`,
		`Arr<Num, Num>`,
		`[2, 3]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[1, 2, 3] =[a, b]`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.ImpossibleMatch),
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
		`{a: 1, b: 2, c: 3} ={a: d, b: e, c: f} d`,
		`Num`,
		`1`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`{a: 1, b: 2, c: 3} ={a: d, b: e, c: f} f`,
		`Num`,
		`3`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`{a: 1, b: 2, c: 3} ={a: d, b: e} e`,
		`Num`,
		`2`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`{a: 1, b: 2, c: 3} =[a, b]`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.ImpossibleMatch),
		),
		t,
	)
	interpreter.TestProgramStr(
		`{a: 1, b: 2, c: 3} ={g: h}`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.ImpossibleMatch),
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
