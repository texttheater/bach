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
		types.NewUnionType(
			types.NullType{},
			types.ObjType{
				Props: map[string]types.Type{
					"start": types.NumType{},
					"0":     types.StrType{},
					"1":     types.NewUnionType(types.NullType{}, types.StrType{}),
					"cs":    types.NewUnionType(types.NullType{}, types.StrType{}),
				},
				Rest: types.VoidType{},
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
		types.NewUnionType(
			types.NullType{},
			types.ObjType{
				Props: map[string]types.Type{
					"start": types.NumType{},
					"0":     types.StrType{},
					"1":     types.NewUnionType(types.NullType{}, types.StrType{}),
					"cs":    types.NewUnionType(types.NullType{}, types.StrType{}),
				},
				Rest: types.VoidType{},
			},
		),
		states.NullValue{},
		nil,
		t,
	)
	interpreter.TestProgram(
		`"abccd" ~^b(?P<cs>*)d~`,
		types.NewUnionType(
			types.NullType{},
			types.ObjType{
				Props: map[string]types.Type{
					"start": types.NumType{},
					"0":     types.StrType{},
					"1":     types.NewUnionType(types.NullType{}, types.StrType{}),
					"cs":    types.NewUnionType(types.NullType{}, types.StrType{}),
				},
				Rest: types.AnyType{},
			},
		),
		nil,
		errors.SyntaxError(
			errors.Code(errors.BadRegexp),
		),
		t,
	)
}
