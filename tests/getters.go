package tests

import (
	"testing"

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
		states.E(
			states.Code(states.NoSuchProperty),
			states.WantType(types.ObjType{
				PropTypeMap: map[string]types.Type{
					"c": types.AnyType{},
				},
				RestType: types.AnyType{},
			}),
			states.GotType(types.ObjType{
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
		states.E(
			states.Code(states.NoSuchIndex),
			states.WantType(&types.NearrType{types.AnyType{}, &types.NearrType{types.AnyType{}, &types.NearrType{types.AnyType{}, &types.NearrType{types.AnyType{}, types.AnyArrType}}}}),
			states.GotType(types.TupType([]types.Type{types.StrType{}, types.StrType{}, types.StrType{}}))),

		t,
	)
	TestProgram(`["a", "b", "c"] @-1`,
		nil,
		nil,
		states.E(
			states.Code(states.BadIndex)),

		t,
	)
	TestProgram(`["a", "b", "c"] @1.5`,
		nil,
		nil,
		states.E(
			states.Code(states.BadIndex)),

		t,
	)
	TestProgram(`"abc" @1`,
		nil,
		nil,
		states.E(
			states.Code(states.NoGetterAllowed)),

		t,
	)
	TestProgram(`24 @1`,
		nil,
		nil,
		states.E(
			states.Code(states.NoGetterAllowed)),

		t,
	)
	TestProgram(`for Any def f Arr<Any> as [] ok f @1`,
		nil,
		nil,
		states.E(
			states.Code(states.NoGetterAllowed)),

		t,
	)
	TestProgram(`for Any def f Arr<Any> as ["a", "b", "c"] ok f @1`,
		nil,
		nil,
		states.E(
			states.Code(states.NoGetterAllowed)),

		t,
	)
}
