package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestRegexp(t *testing.T) {
	TestProgram(
		`"abccd" ~b(?P<cs>c*)d~ must`,
		types.ObjType{
			PropTypeMap: map[string]types.Type{
				"start": types.NumType{},
				"0":     types.StrType{},
				"1":     types.Union(types.NullType{}, types.StrType{}),
				"cs":    types.Union(types.NullType{}, types.StrType{}),
			},
			RestType: types.VoidType{},
		},
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
			types.ObjType{
				PropTypeMap: map[string]types.Type{
					"just": types.ObjType{
						PropTypeMap: map[string]types.Type{
							"start": types.NumType{},
							"0":     types.StrType{},
							"1":     types.Union(types.NullType{}, types.StrType{}),
							"cs":    types.Union(types.NullType{}, types.StrType{}),
						},
						RestType: types.VoidType{},
					},
				},
				RestType: types.AnyType{},
			},
			types.NullType{},
		),
		states.NullValue{},
		nil,
		t,
	)
	TestProgram(
		`"abccd" ~^b(?P<cs>*)d~`,
		types.Union(
			types.ObjType{
				PropTypeMap: map[string]types.Type{
					"just": types.ObjType{
						PropTypeMap: map[string]types.Type{
							"start": types.NumType{},
							"0":     types.StrType{},
							"1":     types.Union(types.NullType{}, types.StrType{}),
							"cs":    types.Union(types.NullType{}, types.StrType{}),
						},
						RestType: types.VoidType{},
					},
				},
				RestType: types.AnyType{},
			},
			types.NullType{},
		),
		nil,
		errors.E(
			errors.Code(errors.BadRegexp),
		),
		t,
	)
	TestProgram(
		`"abcbcdef" findAll~c~`,
		&types.ArrType{
			ElType: types.ObjType{
				PropTypeMap: map[string]types.Type{
					"start": types.NumType{},
					"0":     types.StrType{},
				},
				RestType: types.VoidType{},
			},
		},
		states.NewArrValue([]states.Value{
			states.ObjValue(map[string]*states.Thunk{
				"start": states.ThunkFromValue(states.NumValue(2)),
				"0":     states.ThunkFromValue(states.StrValue("c")),
			}),
			states.ObjValue(map[string]*states.Thunk{
				"start": states.ThunkFromValue(states.NumValue(4)),
				"0":     states.ThunkFromValue(states.StrValue("c")),
			}),
		}),
		nil,
		t,
	)
}
