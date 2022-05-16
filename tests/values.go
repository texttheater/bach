package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
)

func TestValues(t *testing.T) {
	TestProgramStr(
		`false if id then 1 else fatal ok`,
		`Num`,
		``,
		errors.ValueError(
			errors.Code(errors.UnexpectedValue),
			errors.GotValue(states.BoolValue(false)),
		),
		t,
	)
	TestProgramStr(
		`null ==null`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	TestProgramStr(
		`null =={}`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	TestProgramStr(
		`true ==true`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	TestProgramStr(
		`true ==false`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	TestProgramStr(
		`true ==[]`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	TestProgramStr(
		`1 ==1.0`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	TestProgramStr(
		`1 ==2`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	TestProgramStr(
		`57 =="a"`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" =="abc"`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	TestProgramStr(
		`"" =="abc"`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	TestProgramStr(
		`"" ==null`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	TestProgramStr(
		`[false, 1.0, "ab"] ==[false, 1, "a" +"b"]`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	TestProgramStr(
		`[] ==[11]`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	TestProgramStr(
		`["a"] =={a: 1}`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	TestProgramStr(
		`{a: 1, b: 2} =={b: 2, a: 1}`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	TestProgramStr(
		`{a: 1, b: 2} =={a: 2, b: 1}`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	TestProgramStr(
		`{} ==[]`,
		`Bool`,
		`false`,
		nil,
		t,
	)
}
