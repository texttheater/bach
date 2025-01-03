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
		types.Num{},
		states.NumValue(1),
		nil,
		t,
	)
	interpreter.TestProgram(`{a: 1, b: 2} @b`,
		types.Num{},
		states.NumValue(2),
		nil,
		t,
	)
	interpreter.TestProgram(`{a: 1, b: 2} @c`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.NoSuchProperty),
			errors.WantType(types.Obj{
				Props: map[string]types.Type{
					"c": types.Any{},
				},
				Rest: types.Any{},
			}),
			errors.GotType(types.Obj{
				Props: map[string]types.Type{
					"a": types.Num{},
					"b": types.Num{},
				},
				Rest: types.Void{},
			}),
		),
		t,
	)
	interpreter.TestProgram(`["a", "b", "c"] @0`,
		types.Str{},
		states.StrValue("a"),
		nil,
		t,
	)
	interpreter.TestProgram(`["a", "b", "c"] @1`,
		types.Str{},
		states.StrValue("b"),
		nil,
		t,
	)
	interpreter.TestProgram(`["a", "b", "c"] @2`,
		types.Str{},
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
