package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestRegexp(t *testing.T) {
	interpreter.TestProgram(
		`"abccd" ~b(?P<cs>c*)d~`,
		types.NewUnion(
			types.Null{},
			types.Obj{
				Props: map[string]types.Type{
					"start": types.Num{},
					"0":     types.Str{},
					"1":     types.NewUnion(types.Null{}, types.Str{}),
					"cs":    types.NewUnion(types.Null{}, types.Str{}),
				},
				Rest: types.Void{},
			},
		),
		states.ObjValueFromMap(map[string]states.Value{
			"start": states.NumValue(1),
			"0":     states.StrValue("bccd"),
			"1":     states.StrValue("cc"),
			"cs":    states.StrValue("cc"),
		}),
		nil,
		t,
	)
	interpreter.TestProgram(
		`"def" ~^b(?P<cs>c*)d~`,
		types.NewUnion(
			types.Null{},
			types.Obj{
				Props: map[string]types.Type{
					"start": types.Num{},
					"0":     types.Str{},
					"1":     types.NewUnion(types.Null{}, types.Str{}),
					"cs":    types.NewUnion(types.Null{}, types.Str{}),
				},
				Rest: types.Void{},
			},
		),
		states.NullValue{},
		nil,
		t,
	)
	interpreter.TestProgram(
		`"abccd" ~^b(?P<cs>*)d~`,
		types.NewUnion(
			types.Null{},
			types.Obj{
				Props: map[string]types.Type{
					"start": types.Num{},
					"0":     types.Str{},
					"1":     types.NewUnion(types.Null{}, types.Str{}),
					"cs":    types.NewUnion(types.Null{}, types.Str{}),
				},
				Rest: types.Any{},
			},
		),
		nil,
		errors.SyntaxError(
			errors.Code(errors.BadRegexp),
		),
		t,
	)
}
