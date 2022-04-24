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
	TestProgram(
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
	TestProgram(
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
	TestProgramStr(
		`"zabacad" split~a~`,
		`Arr<Str>`,
		`["z", "b", "c", "d"]`,
		nil,
		t,
	)
	TestProgramStr(
		`"zabaca" split~a~`,
		`Arr<Str>`,
		`["z", "b", "c", ""]`,
		nil,
		t,
	)
	TestProgramStr(
		`"abacad" split~a~`,
		`Arr<Str>`,
		`["", "b", "c", "d"]`,
		nil,
		t,
	)
	TestProgramStr(
		`"abaca" split~a~`,
		`Arr<Str>`,
		`["", "b", "c", ""]`,
		nil,
		t,
	)
	TestProgramStr(
		`"abaca" split~~`,
		`Arr<Str>`,
		`["a", "b", "a", "c", "a"]`,
		nil,
		t,
	)
	TestProgramStr(
		`"你好" split~~`,
		`Arr<Str>`,
		`["你", "好"]`,
		nil,
		t,
	)
	TestProgramStr(
		`"" split~a~`,
		`Arr<Str>`,
		`[""]`,
		nil,
		t,
	)
	TestProgramStr(
		`"" split~~`,
		`Arr<Str>`,
		`[]`,
		nil,
		t,
	)
}
