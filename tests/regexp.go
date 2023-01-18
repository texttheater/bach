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
		`"zabacad" split(~a~, 1)`,
		`Arr<Str>`,
		`["z", "bacad"]`,
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
		`"zabaca" split(~a~, 1)`,
		`Arr<Str>`,
		`["z", "baca"]`,
		nil,
		t,
	)
	TestProgramStr(
		`"zabaca" split(~a~, 3)`,
		`Arr<Str>`,
		`["z", "b", "c", ""]`,
		nil,
		t,
	)
	TestProgramStr(
		`"zabaca" split(~a~, 4)`,
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
		`"abacad" split(~a~, 1)`,
		`Arr<Str>`,
		`["", "bacad"]`,
		nil,
		t,
	)
	TestProgramStr(
		`"abacad" split(~a~, 2)`,
		`Arr<Str>`,
		`["", "b", "cad"]`,
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
		`"abaca" split(~a~, 1)`,
		`Arr<Str>`,
		`["", "baca"]`,
		nil,
		t,
	)
	TestProgramStr(
		`"abaca" split(~a~, 2)`,
		`Arr<Str>`,
		`["", "b", "ca"]`,
		nil,
		t,
	)
	TestProgramStr(
		`"abaca" split(~a~, 3)`,
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
		`"abaca" split(~~, 2)`,
		`Arr<Str>`,
		`["a", "b", "aca"]`,
		nil,
		t,
	)
	TestProgramStr(
		`"abaca" split(~~, 1000)`,
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
		`"你好" split(~~, 0)`,
		`Arr<Str>`,
		`["你好"]`,
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
		`"" split(~a~, 0)`,
		`Arr<Str>`,
		`[""]`,
		nil,
		t,
	)
	TestProgramStr(
		`"" split(~a~, 1)`,
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
	TestProgramStr(
		`"" split(~~, 0)`,
		`Arr<Str>`,
		`[]`,
		nil,
		t,
	)
	TestProgramStr(
		`"" split(~~, 1)`,
		`Arr<Str>`,
		`[]`,
		nil,
		t,
	)
	TestProgramStr(
		`"" replaceFirst(~bc~, "hurz")`,
		`Str`,
		`""`,
		nil,
		t,
	)
	TestProgramStr(
		`"abcd" replaceFirst(~bc~, "hurz")`,
		`Str`,
		`"ahurzd"`,
		nil,
		t,
	)
	TestProgramStr(
		`"abcdabcd" replaceFirst(~bc~, "hurz")`,
		`Str`,
		`"ahurzdabcd"`,
		nil,
		t,
	)
	TestProgramStr(
		`"abcd" replaceFirst(~bd~, "hurz")`,
		`Str`,
		`"abcd"`,
		nil,
		t,
	)
	TestProgramStr(
		`"" replaceAll(~a~, "b")`,
		`Str`,
		`""`,
		nil,
		t,
	)
	TestProgramStr(
		`"ababa" replaceAll(~a~, "b")`,
		`Str`,
		`"bbbbb"`,
		nil,
		t,
	)
	TestProgramStr(
		`"babab" replaceAll(~a~, "b")`,
		`Str`,
		`"bbbbb"`,
		nil,
		t,
	)
}
