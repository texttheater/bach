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
		`if true then {a: 1} else {b: 2} ok is {a: Num} then true elseIs {a: Num, b: Num} then false ok`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.ImpossibleMatch),
		),
		t,
	)
	TestProgram(
		`if true then {a: 1} else {b: 2} ok is {a: Num a} then a elseIs {b: Num b} then b ok`,
		types.NumType{},
		values.NumValue(1),
		nil,
		t,
	)
	TestProgram(
		`if true then {a: 1} else {b: "s"} ok is {a: Num a} then a elseIs {b: Num b} then b ok`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.ImpossibleMatch),
		),
		t,
	)
	TestProgram(
		`if true then {a: 1} elseIf true then {b: "s"} else {c: true} ok is {a: Num a} then a elseIs {b: Num b} then b ok`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.NonExhaustiveMatch),
		),
		t,
	)
	TestProgram(
		`if true then {a: 1} else {b: 2} ok`,
		types.Union(types.NewObjType(map[string]types.Type{"a": types.NumType{}}), types.NewObjType(map[string]types.Type{"b": types.NumType{}})),
		values.ObjValue{"a": values.NumValue(1)},
		nil,
		t,
	)
	// TODO In the following program, the type of x could be inferred in
	// each pattern but isn't. Is this a problem? It is because our object
	// types currently cannot represent the *absence* of attributes, so any
	// attribute not present in the type has Any type by default.
	TestProgram(
		`if true then {a: 1} else {b: 2} ok is {a: Num x} with x ==1 then true elseIf false then false elseIs {a: Num x} then x elseIs {b: Num x} then x ok`,
		types.Union(types.BoolType{}, types.NumType{}),
		values.BoolValue(true),
		nil,
		t,
	)
}
