package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestObjects(t *testing.T) {
	TestProgram(
		`{} get"a"`,
		types.Union(
			types.ObjType{
				PropTypeMap: map[string]types.Type{
					"just": types.VoidType{},
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
		`{a: 1} get"a"`,
		types.Union(
			types.ObjType{
				PropTypeMap: map[string]types.Type{
					"just": types.NumType{},
				},
				RestType: types.AnyType{},
			},
			types.NullType{},
		),
		states.ObjValue(map[string]*states.Thunk{
			"just": states.ThunkFromValue(states.NumValue(1)),
		}),
		nil,
		t,
	)
	TestProgram(
		`{a: 1, b: "hey"} get"a"`,
		types.Union(
			types.ObjType{
				PropTypeMap: map[string]types.Type{
					"just": types.Union(
						types.NumType{},
						types.StrType{},
					),
				},
				RestType: types.AnyType{},
			},
			types.NullType{},
		),
		states.ObjValue(map[string]*states.Thunk{
			"just": states.ThunkFromValue(states.NumValue(1)),
		}),
		nil,
		t,
	)
	TestProgram(
		`{a: 1, b: "hey", c: false} get"a"`,
		types.Union(
			types.ObjType{
				PropTypeMap: map[string]types.Type{
					"just": types.Union(
						types.NumType{},
						types.Union(
							types.StrType{},
							types.BoolType{},
						),
					),
				},
				RestType: types.AnyType{},
			},
			types.NullType{},
		),
		states.ObjValue(map[string]*states.Thunk{
			"just": states.ThunkFromValue(states.NumValue(1)),
		}),
		nil,
		t,
	)
	TestProgram(
		`{1: "a"} get(1)`,
		types.Union(
			types.ObjType{
				PropTypeMap: map[string]types.Type{
					"just": types.StrType{},
				},
				RestType: types.AnyType{},
			},
			types.NullType{},
		),
		states.ObjValue(map[string]*states.Thunk{
			"just": states.ThunkFromValue(states.StrValue("a")),
		}),
		nil,
		t,
	)
	TestProgram(
		`{1.5: "a"} get(1.5)`,
		types.Union(
			types.ObjType{
				PropTypeMap: map[string]types.Type{
					"just": types.StrType{},
				},
				RestType: types.AnyType{},
			},
			types.NullType{},
		),
		states.ObjValue(map[string]*states.Thunk{
			"just": states.ThunkFromValue(states.StrValue("a")),
		}),
		nil,
		t,
	)
}
