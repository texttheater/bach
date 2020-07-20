package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestMatchingType(t *testing.T) {
	TestProgram(
		`if true then 2 else "two" ok is Num then true else false ok`,
		types.BoolType{},
		states.BoolValue(true),
		nil,
		t,
	)
	//TestProgram(
	//	`if true then 2 else "two" ok is Num then true ok`,
	//	nil,
	//	nil,
	//	errors.E(
	//		errors.Code(errors.NonExhaustiveMatch),
	//		errors.WantType(types.VoidType{}),
	//		errors.GotType(types.StrType{}),
	//	),
	//	t,
	//)
	TestProgram(
		`if true then 2 else "two" ok is Num then true elis Str then false ok`,
		types.BoolType{},
		states.BoolValue(true),
		nil,
		t,
	)
	TestProgram(
		`if true then 2 else "two" ok is Num then true elis Str then false else false ok`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.UnreachableElseClause)),

		t,
	)
	TestProgram(
		`[1, 2, 3] is Arr<Num> then each +1 all ok`,
		&types.ArrType{types.NumType{}},
		states.NewArrValue([]states.Value{states.NumValue(2), states.NumValue(3), states.NumValue(4)}),
		nil,
		t,
	)
	TestProgram( // Intersective matching: pattern says Arr<Any> but Bach knows it got Arr<Num>
		`[1, 2, 3] is Arr<Any> then each +1 all ok`,
		&types.ArrType{types.NumType{}},
		states.NewArrValue([]states.Value{states.NumValue(2), states.NumValue(3), states.NumValue(4)}),
		nil,
		t,
	)
}
