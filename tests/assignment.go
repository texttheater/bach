package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
)

func TestAssignment(t *testing.T) {
	TestProgramStr(
		`1 +1 =a 3 *2 +a`,
		`Num`,
		`8`,
		nil,
		t,
	)
	TestProgramStr(
		`1 +1 ==2 =p 1 +1 ==1 =q p ==q not`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] =[a, b, c] a`,
		`Num`,
		`1`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] =[a, b, c] c`,
		`Num`,
		`3`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] =[a;r] r`,
		`Tup<Num, Num>`,
		`[2, 3]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] =[a, b]`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.ImpossibleMatch),
		),
		t,
	)
	TestProgramStr(
		`for Arr<Num> def f Num as =[a, b] a ok`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.NonExhaustiveMatch),
		),
		t,
	)
}
