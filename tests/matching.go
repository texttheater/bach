package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestMatching(t *testing.T) {
	TestProgram(
		`if true then 2 else "two" ok is Num then true else false ok`,
		types.BoolType{},
		values.BoolValue(true),
		nil,
		t)
	TestProgram(
		`if true then 2 else "two" ok is Num then true ok`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.NonExhaustiveMatch),
			errors.WantType(types.VoidType{}),
			errors.GotType(types.StrType{}),
		),
		t)
	TestProgram(
		`if true then 2 else "two" ok is Num then true elis Str then false ok`,
		types.BoolType{},
		values.BoolValue(true),
		nil,
		t)
	TestProgram(
		`if true then 2 else "two" ok is Num then true elis Str then false else false ok`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.UnreachableElseClause),
		),
		t)
	TestProgram(
		`[1, 2, 3] is Seq<Num> then each +1 all arr ok`,
		&types.ArrType{types.NumType{}},
		values.ArrValue([]values.Value{values.NumValue(2), values.NumValue(3), values.NumValue(4)}),
		nil,
		t)
	TestProgram( // Intersective matching: pattern says Seq<Any> but Bach knows it got Seq<Num>
		`[1, 2, 3] is Seq<Any> then each +1 all arr ok`,
		&types.ArrType{types.NumType{}},
		values.ArrValue([]values.Value{values.NumValue(2), values.NumValue(3), values.NumValue(4)}),
		nil,
		t)
}
