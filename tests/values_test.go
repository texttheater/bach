package tests_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/tests"
)

func TestValues(t *testing.T) {
	tests.TestProgramStr(
		`false if id then 1 else fatal ok`,
		`Num`,
		``,
		errors.ValueError(
			errors.Code(errors.UnexpectedValue),
			errors.GotValue(states.BoolValue(false)),
		),
		t,
	)
	tests.TestProgramStr(
		`null ==null`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`null =={}`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`true ==true`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`true ==false`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`true ==[]`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`1 ==1.0`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`1 ==2`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`57 =="a"`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" =="abc"`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"" =="abc"`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"" ==null`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`[false, 1.0, "ab"] ==[false, 1, "a" +"b"]`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`[] ==[11]`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`["a"] =={a: 1}`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`{a: 1, b: 2} =={b: 2, a: 1}`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`{a: 1, b: 2} =={a: 2, b: 1}`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`{} ==[]`,
		`Bool`,
		`false`,
		nil,
		t,
	)
}
