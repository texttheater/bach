package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestMatchingObj(t *testing.T) {
	TestProgram(
		`{a: 1} is {a: Num a} then a ok`,
		types.NumType{},
		states.NumValue(1),
		nil,
		t,
	)
	TestProgram(
		`{a: 1, b: 2} is {a: Num a} then a ok`,
		types.NumType{},
		states.NumValue(1),
		nil,
		t,
	)
	TestProgram(
		`if true then {a: 1} else {b: 2} ok is {a: Num} then true elis {a: Num, b: Num} then false ok`,
		nil,
		nil,
		states.E(
			states.Code(states.ImpossibleMatch)),

		t,
	)
	TestProgram(
		`if true then {a: 1} else {b: 2} ok is {a: Num a} then a elis {b: Num b} then b ok`,
		types.NumType{},
		states.NumValue(1),
		nil,
		t,
	)
	TestProgram(
		`if true then {a: 1} else {b: "s"} ok is {a: Num a} then a elis {b: Num b} then b ok`,
		nil,
		nil,
		states.E(
			states.Code(states.ImpossibleMatch)),

		t,
	)
	//TestProgram(
	//	`if true then {a: 1} elif true then {b: "s"} else {c: true} ok is {a: Num a} then a elis {b: Num b} then b ok`,
	//	nil,
	//	nil,
	//	errors.E(
	//		errors.Code(errors.NonExhaustiveMatch),
	//	),
	//	t,
	//)
	TestProgram(
		`if true then {a: 1} else {b: 2} ok`,
		types.Union(
			types.ObjType{
				PropTypeMap: map[string]types.Type{
					"a": types.NumType{},
				},
				RestType: types.VoidType{},
			},
			types.ObjType{
				PropTypeMap: map[string]types.Type{
					"b": types.NumType{},
				},
				RestType: types.VoidType{},
			},
		),
		states.ObjValueFromMap(map[string]states.Value{
			"a": states.NumValue(1),
		}),
		nil,
		t,
	)
	// TODO In the following program, the type of x could be inferred in
	// each pattern but isn't. Is this a problem? It is because our object
	// types currently cannot represent the *absence* of attributes, so any
	// attribute not present in the type has Any type by default.
	TestProgram(
		`if true then {a: 1} else {b: 2} ok is {a: Num x} with x ==1 then true elif false then false elis {a: Num x} then x elis {b: Num x} then x ok`,
		types.Union(types.BoolType{}, types.NumType{}),
		states.BoolValue(true),
		nil,
		t,
	)
}
