package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestGetters(t *testing.T) {
	interpreter.TestProgram(`{a: 1, b: 2} @a`,
		types.NumType{},
		states.NumValue(1),
		nil,
		t,
	)
	interpreter.TestProgram(`{a: 1, b: 2} @b`,
		types.NumType{},
		states.NumValue(2),
		nil,
		t,
	)
	interpreter.TestProgram(`{a: 1, b: 2} @c`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.NoSuchProperty),
			errors.WantType(types.ObjType{
				Props: map[string]types.Type{
					"c": types.AnyType{},
				},
				Rest: types.AnyType{},
			}),
			errors.GotType(types.ObjType{
				Props: map[string]types.Type{
					"a": types.NumType{},
					"b": types.NumType{},
				},
				Rest: types.VoidType{},
			}),
		),
		t,
	)
	interpreter.TestProgram(`["a", "b", "c"] @0`,
		types.StrType{},
		states.StrValue("a"),
		nil,
		t,
	)
	interpreter.TestProgram(`["a", "b", "c"] @1`,
		types.StrType{},
		states.StrValue("b"),
		nil,
		t,
	)
	interpreter.TestProgram(`["a", "b", "c"] @2`,
		types.StrType{},
		states.StrValue("c"),
		nil,
		t,
	)
	interpreter.TestProgram(`["a", "b", "c"] @3`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.NoSuchIndex),
		),
		t,
	)
	interpreter.TestProgramStr(
		`["a", "b", "c"] @-1`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.BadIndex),
		),
		t,
	)
	interpreter.TestProgram(`["a", "b", "c"] @1.5`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.BadIndex),
		),
		t,
	)
	interpreter.TestProgram(`"abc" @1`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.NoGetterAllowed),
		),
		t,
	)
	interpreter.TestProgram(`24 @1`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.NoGetterAllowed),
		),
		t,
	)
	interpreter.TestProgram(`for Any def f Arr<Any...> as [] ok f @1`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.NoGetterAllowed),
		),
		t,
	)
	interpreter.TestProgram(`for Any def f Arr<Any...> as ["a", "b", "c"] ok f @1`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.NoGetterAllowed),
		),
		t,
	)
}
