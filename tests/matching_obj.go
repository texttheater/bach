package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestMatchingObj(t *testing.T) {
	TestProgram(
		`{a: 1} is {a: Num a} then a ok`,
		types.NumType{},
		values.NumValue(1),
		nil,
		t,
	)
	TestProgram(
		`{a: 1, b: 2} is {a: Num a} then a ok`,
		types.NumType{},
		values.NumValue(1),
		nil,
		t,
	)
	TestProgram(
		`if true then {a: 1} else {b: 2} ok is {a: Num} then true elis {a: Num, b: Num} then false ok`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.ImpossibleMatch),
		),
		t,
	)
	TestProgram(
		`if true then {a: 1} else {b: 2} ok is {a: Num a} then a elis {b: Num b} then b ok`,
		types.NumType{},
		values.NumValue(1),
		nil,
		t,
	)
	TestProgram(
		`if true then {a: 1} else {b: "s"} ok is {a: Num a} then a elis {b: Num b} then b ok`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.ImpossibleMatch),
		),
		t,
	)
	TestProgram(
		`if true then {a: 1} elif true then {b: "s"} else {c: true} ok is {a: Num a} then a elis {b: Num b} then b ok`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.NonExhaustiveMatch),
		),
		t,
	)
}
