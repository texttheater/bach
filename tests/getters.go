package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestGetters(t *testing.T) {
	TestProgram(`{a: 1, b: 2} @a`,
		types.NumType{},
		states.NumValue(1),
		nil,
		t,
	)
	TestProgram(`{a: 1, b: 2} @b`,
		types.NumType{},
		states.NumValue(2),
		nil,
		t,
	)
	TestProgram(`{a: 1, b: 2} @c`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.NoSuchProperty),
			errors.WantType(types.ObjType{
				PropTypeMap: map[string]types.Type{
					"c": types.AnyType{},
				},
				RestType: types.AnyType{},
			}),
			errors.GotType(types.ObjType{
				PropTypeMap: map[string]types.Type{
					"a": types.NumType{},
					"b": types.NumType{},
				},
				RestType: types.VoidType{},
			})),

		t,
	)
	TestProgram(`["a", "b", "c"] @0`,
		types.StrType{},
		states.StrValue("a"),
		nil,
		t,
	)
	TestProgram(`["a", "b", "c"] @1`,
		types.StrType{},
		states.StrValue("b"),
		nil,
		t,
	)
	TestProgram(`["a", "b", "c"] @2`,
		types.StrType{},
		states.StrValue("c"),
		nil,
		t,
	)
	TestProgram(`["a", "b", "c"] @3`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.NoSuchIndex),
		),
		t,
	)
	TestProgramStr(
		`["a", "b", "c"] @-1`,
		`Str`,
		`"c"`,
		nil,
		t,
	)
	TestProgramStr(
		`["a", "b", "c"] @-2`,
		`Str`,
		`"b"`,
		nil,
		t,
	)
	TestProgramStr(
		`["a", "b", "c"] @-3`,
		`Str`,
		`"a"`,
		nil,
		t,
	)
	TestProgram(`["a", "b", "c"] @1.5`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.BadIndex)),

		t,
	)
	TestProgram(`"abc" @1`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.NoGetterAllowed)),

		t,
	)
	TestProgram(`24 @1`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.NoGetterAllowed)),

		t,
	)
	TestProgram(`for Any def f Arr<Any> as [] ok f @1`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.NoGetterAllowed)),

		t,
	)
	TestProgram(`for Any def f Arr<Any> as ["a", "b", "c"] ok f @1`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.NoGetterAllowed)),

		t,
	)
}
