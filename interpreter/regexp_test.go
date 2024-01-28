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
	interpreter.TestProgramStr(
		`"zabacad" reSplit~a~`,
		`Arr<Str>`,
		`["z", "b", "c", "d"]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"zabacad" reSplit(~a~, 1)`,
		`Arr<Str>`,
		`["z", "bacad"]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"zabaca" reSplit~a~`,
		`Arr<Str>`,
		`["z", "b", "c", ""]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"zabaca" reSplit(~a~, 1)`,
		`Arr<Str>`,
		`["z", "baca"]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"zabaca" reSplit(~a~, 3)`,
		`Arr<Str>`,
		`["z", "b", "c", ""]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"zabaca" reSplit(~a~, 4)`,
		`Arr<Str>`,
		`["z", "b", "c", ""]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abacad" reSplit~a~`,
		`Arr<Str>`,
		`["", "b", "c", "d"]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abacad" reSplit(~a~, 1)`,
		`Arr<Str>`,
		`["", "bacad"]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abacad" reSplit(~a~, 2)`,
		`Arr<Str>`,
		`["", "b", "cad"]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abaca" reSplit~a~`,
		`Arr<Str>`,
		`["", "b", "c", ""]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abaca" reSplit(~a~, 1)`,
		`Arr<Str>`,
		`["", "baca"]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abaca" reSplit(~a~, 2)`,
		`Arr<Str>`,
		`["", "b", "ca"]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abaca" reSplit(~a~, 3)`,
		`Arr<Str>`,
		`["", "b", "c", ""]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abaca" reSplit~~`,
		`Arr<Str>`,
		`["a", "b", "a", "c", "a"]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abaca" reSplit(~~, 2)`,
		`Arr<Str>`,
		`["a", "b", "aca"]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abaca" reSplit(~~, 1000)`,
		`Arr<Str>`,
		`["a", "b", "a", "c", "a"]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"你好" reSplit~~`,
		`Arr<Str>`,
		`["你", "好"]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"你好" reSplit(~~, 0)`,
		`Arr<Str>`,
		`["你好"]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"" reSplit~a~`,
		`Arr<Str>`,
		`[""]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"" reSplit(~a~, 0)`,
		`Arr<Str>`,
		`[""]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"" reSplit(~a~, 1)`,
		`Arr<Str>`,
		`[""]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"" reSplit~~`,
		`Arr<Str>`,
		`[]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"" reSplit(~~, 0)`,
		`Arr<Str>`,
		`[]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"" reSplit(~~, 1)`,
		`Arr<Str>`,
		`[]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"" reReplaceFirst(~bc~, "hurz")`,
		`Str`,
		`""`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abcd" reReplaceFirst(~bc~, "hurz")`,
		`Str`,
		`"ahurzd"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abcdabcd" reReplaceFirst(~bc~, "hurz")`,
		`Str`,
		`"ahurzdabcd"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abcd" reReplaceFirst(~bd~, "hurz")`,
		`Str`,
		`"abcd"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"" reReplaceAll(~a~, "b")`,
		`Str`,
		`""`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"ababa" reReplaceAll(~a~, "b")`,
		`Str`,
		`"bbbbb"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"babab" reReplaceAll(~a~, "b")`,
		`Str`,
		`"bbbbb"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`" a b c " reReplaceFirst(~[abc]~, "({@0})")`,
		`Str`,
		`" (a) b c "`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`" a b c " reReplaceAll(~[abc]~, "({@0})")`,
		`Str`,
		`" (a) (b) (c) "`,
		nil,
		t,
	)
}
