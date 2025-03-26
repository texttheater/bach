package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestMatchingArr(t *testing.T) {
	interpreter.TestProgram(
		`[1, 2, 3] is [Num, Num, Num] then true ok`,
		types.BoolType{},
		states.BoolValue(true),
		nil,
		t,
	)
	interpreter.TestProgram(
		`[1, 2, 3] is [Num, Num, Num] then true else false ok`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.UnreachableElseClause),
		),
		t,
	)
	interpreter.TestProgram(
		`[1, 2, 3] each(id) is [Num, Num, Num] then true else false ok`,
		types.BoolType{},
		states.BoolValue(true),
		nil,
		t,
	)
	interpreter.TestProgram(
		`[1, "a"] is [Num, Str] then true ok`,
		types.BoolType{},
		states.BoolValue(true),
		nil,
		t,
	)
	interpreter.TestProgram(
		`[1, "a"] is [Num a, Str b] then a ok`,
		types.NumType{},
		states.NumValue(1),
		nil,
		t,
	)
	interpreter.TestProgram(
		`[1, "a"] is [Num a, Str b] then b ok`,
		types.StrType{},
		states.StrValue("a"),
		nil,
		t,
	)
	interpreter.TestProgram(
		`[[1]] is [[Any x]] then x ok`,
		types.NumType{},
		states.NumValue(1),
		nil,
		t,
	)
	interpreter.TestProgram(
		`[[1]] is [[x]] then x ok`,
		types.NumType{},
		states.NumValue(1),
		nil,
		t,
	)
	interpreter.TestProgram(
		`if true then [1] else [2] ok is [Num a] then a ok`,
		types.NumType{},
		states.NumValue(1),
		nil,
		t,
	)
	//interpreter.TestProgram(
	//	`if true then [1] else ["2"] ok is [Num a] then a ok`,
	//	nil,
	//	nil,
	//	errors.TypeError(
	//		errors.Code(errors.NonExhaustiveMatch),
	//	),
	//	t,
	//)
	interpreter.TestProgram(
		`if true then [1] else ["2"] ok is [Num a] then a elis [Str a] then a ok`,
		types.NewUnionType(types.NumType{}, types.StrType{}),
		states.NumValue(1),
		nil,
		t,
	)
	interpreter.TestProgram(
		`[] is [a] then a ok`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.ImpossibleMatch),
		),
		t,
	)
	interpreter.TestProgram(
		`[1] is [a, b] then a ok`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.ImpossibleMatch),
		),
		t,
	)
	interpreter.TestProgram(
		`[1, 2, 3] is [head;tail] then tail ok`,
		&types.NearrType{
			Head: types.NumType{},
			Tail: &types.NearrType{
				Head: types.NumType{},
				Tail: types.VoidArrType,
			},
		},
		states.NewArrValue(
			[]states.Value{
				states.NumValue(2),
				states.NumValue(3),
			},
		),
		nil,
		t,
	)
	interpreter.TestProgram(
		`[1, 2, 3] is [a, b;rest] then rest ok`,
		&types.NearrType{
			Head: types.NumType{},
			Tail: types.VoidArrType,
		},
		states.NewArrValue(
			[]states.Value{
				states.NumValue(3),
			},
		),
		nil,
		t,
	)
	interpreter.TestProgram(
		`[1, 2, 3] is [a, b, c;rest] then rest ok`,
		types.VoidArrType,
		states.NewArrValue(
			[]states.Value{},
		),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Arr<Num...> def plusOne Arr<Num...> as is [head;tail] then [head +1;tail plusOne] else [] ok ok [1, 2] plusOne`,
		&types.ArrType{
			El: types.NumType{},
		},
		states.NewArrValue([]states.Value{
			states.NumValue(2),
			states.NumValue(3),
		}),
		nil,
		t,
	)
}
