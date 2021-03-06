package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestRegexp(t *testing.T) {
	TestProgram(
		`"abccd" ~b(?P<cs>c*)d~`,
		types.Union(
			types.NullType{},
			types.ObjType{
				PropTypeMap: map[string]types.Type{
					"start": types.NumType{},
					"0":     types.StrType{},
					"1":     types.Union(types.NullType{}, types.StrType{}),
					"cs":    types.Union(types.NullType{}, types.StrType{}),
				},
				RestType: types.VoidType{},
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
	TestProgram(
		`"def" ~^b(?P<cs>c*)d~`,
		types.Union(
			types.NullType{},
			types.ObjType{
				PropTypeMap: map[string]types.Type{
					"start": types.NumType{},
					"0":     types.StrType{},
					"1":     types.Union(types.NullType{}, types.StrType{}),
					"cs":    types.Union(types.NullType{}, types.StrType{}),
				},
				RestType: types.VoidType{},
			},
		),
		states.NullValue{},
		nil,
		t,
	)
	TestProgram(
		`"abccd" ~^b(?P<cs>*)d~`,
		types.Union(
			types.NullType{},
			types.ObjType{
				PropTypeMap: map[string]types.Type{
					"start": types.NumType{},
					"0":     types.StrType{},
					"1":     types.Union(types.NullType{}, types.StrType{}),
					"cs":    types.Union(types.NullType{}, types.StrType{}),
				},
				RestType: types.AnyType{},
			},
		),
		nil,
		errors.E(
			errors.Code(errors.BadRegexp)),

		t,
	)
}
