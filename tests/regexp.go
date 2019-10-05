package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestRegexp(t *testing.T) {
	TestProgram(
		`"abccd" ~b(?P<cs>c*)d~`,
		types.Union(
			types.NullType{},
			types.NewObjType(
				map[string]types.Type{
					"start": types.NumType{},
					"0":     types.StrType{},
					"1":     types.Union(types.NullType{}, types.StrType{}),
					"cs":    types.Union(types.NullType{}, types.StrType{}),
				},
			),
		),
		values.ObjValue(
			map[string]values.Value{
				"start": values.NumValue(1),
				"0":     values.StrValue("bccd"),
				"1":     values.StrValue("cc"),
				"cs":    values.StrValue("cc"),
			},
		),
		nil,
		t,
	)
	TestProgram(
		`"def" ~^b(?P<cs>c*)d~`,
		types.Union(
			types.NullType{},
			types.NewObjType(
				map[string]types.Type{
					"start": types.NumType{},
					"0":     types.StrType{},
					"1":     types.Union(types.NullType{}, types.StrType{}),
					"cs":    types.Union(types.NullType{}, types.StrType{}),
				},
			),
		),
		values.NullValue{},
		nil,
		t,
	)
	TestProgram(
		`"abccd" ~^b(?P<cs>*)d~`,
		types.Union(
			types.NullType{},
			types.NewObjType(
				map[string]types.Type{
					"start": types.NumType{},
					"0":     types.StrType{},
					"1":     types.Union(types.NullType{}, types.StrType{}),
					"cs":    types.Union(types.NullType{}, types.StrType{}),
				},
			),
		),
		nil,
		errors.E(
			errors.Code(errors.BadRegexp),
		),
		t,
	)
}
