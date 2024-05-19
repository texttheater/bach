package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestArrays(t *testing.T) {
	interpreter.TestProgram(
		`[]`,
		types.NewTup([]types.Type{}),
		states.NewArrValue([]states.Value{}),
		nil,
		t,
	)
	interpreter.TestProgram(
		`[1]`,
		types.NewTup([]types.Type{types.Num{}}),
		states.NewArrValue([]states.Value{states.NumValue(1)}),
		nil,
		t,
	)
	interpreter.TestProgram(
		`[1, 2, 3]`,
		types.NewTup([]types.Type{types.Num{}, types.Num{}, types.Num{}}),
		states.NewArrValue([]states.Value{states.NumValue(1),
			states.NumValue(2),
			states.NumValue(3)}),
		nil,
		t,
	)
	interpreter.TestProgram(
		`[1, "a"]`,
		types.NewTup([]types.Type{types.Num{}, types.Str{}}),
		states.NewArrValue([]states.Value{states.NumValue(1),
			states.StrValue("a")}),
		nil,
		t,
	)
	interpreter.TestProgram(
		`[[1, 2], ["a", "b"]]`,
		types.NewTup([]types.Type{types.NewTup([]types.Type{types.Num{}, types.Num{}}), types.NewTup([]types.Type{types.Str{}, types.Str{}})}),
		states.NewArrValue([]states.Value{states.NewArrValue([]states.Value{states.NumValue(1),
			states.NumValue(2)}),
			states.NewArrValue([]states.Value{states.StrValue("a"), states.StrValue("b")})}),
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`for Arr<Num> def f Arr<Num> as =x x ok [1, 2, 3] f`,
		`Arr<Num>`,
		`[1, 2, 3]`,
		nil,
		t,
	)
	interpreter.TestProgram(
		`1 each(*2)`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.NoSuchFunction),
			errors.InputType(types.Num{}),
			errors.Name("each"),
			errors.NumParams(1),
		),
		t,
	)
	interpreter.TestProgram(
		`[1;[]]`,
		types.NewTup(
			[]types.Type{
				types.Num{},
			},
		),
		states.NewArrValue(
			[]states.Value{
				states.NumValue(1),
			},
		),
		nil,
		t,
	)
	interpreter.TestProgram(
		`[3, 4] =rest [1, 2;rest]`,
		types.NewTup(
			[]types.Type{
				types.Num{},
				types.Num{},
				types.Num{},
				types.Num{},
			},
		),
		states.NewArrValue(
			[]states.Value{
				states.NumValue(1),
				states.NumValue(2),
				states.NumValue(3),
				states.NumValue(4),
			},
		),
		nil,
		t,
	)
	interpreter.TestProgram(
		`[1;2]`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.RestRequiresArrType),
			errors.WantType(types.AnyArr),
			errors.GotType(types.Num{}),
		),
		t,
	)
	interpreter.TestProgram(
		`[1 if ==2 then true else fatal ok] out`,
		nil,
		nil,
		errors.ValueError(
			errors.Code(errors.UnexpectedValue),
			errors.GotValue(states.NumValue(1)),
		),
		t,
	)
	interpreter.TestProgram(
		`[true, 1 if ==2 then true else fatal ok] out`,
		nil,
		nil,
		errors.ValueError(
			errors.Code(errors.UnexpectedValue),
			errors.GotValue(states.NumValue(1)),
		),
		t,
	)
	interpreter.TestProgramStr(
		`[1, 2, 3] take(2)`,
		`Arr<Num>`,
		`[1, 2]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[1, 2, 3] take(1)`,
		`Arr<Num>`,
		`[1]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[1, 2, 3] take(0)`,
		`Arr<Num>`,
		`[]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[1, 2, 3] take(-1)`,
		`Arr<Num>`,
		`[]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[1, 2, 3] take(4)`,
		`Arr<Num>`,
		`[1, 2, 3]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[1, 2, 3] take(3)`,
		`Arr<Num>`,
		`[1, 2, 3]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[1, 2, 3] rev`,
		`Arr<Num>`,
		`[3, 2, 1]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[] rev`,
		`Tup<>`,
		`[]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[7, 3, 2, 5, 2] sort`,
		`Arr<Num>`,
		`[2, 2, 3, 5, 7]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"Zwölf Boxkämpfer jagen Victor quer über den großen Sylter Deich . Voilà !" fields sort`,
		`Arr<Str>`,
		`["!", ".", "Boxkämpfer", "Deich", "Sylter", "Victor", "Voilà", "Zwölf", "den", "großen", "jagen", "quer", "über"]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"Zwölf Boxkämpfer jagen Victor quer über den großen Sylter Deich . Voilà !" fields sort(>)`,
		`Arr<Str>`,
		`["über", "quer", "jagen", "großen", "den", "Zwölf", "Voilà", "Victor", "Sylter", "Deich", "Boxkämpfer", ".", "!"]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[7, 3, 2, 5] sort(>)`,
		`Arr<Num>`,
		`[7, 5, 3, 2]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[{a: 7}, {a: 3}, {a: 2}, {a: 5}] for Obj<a: Num, Void> def <(other Obj<a: Num, Void>) Bool as @a <(other @a) ok sort(<)`,
		`Arr<Obj<a: Num, Void>>`,
		`[{a: 2}, {a: 3}, {a: 5}, {a: 7}]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[{a: 7, b: 2}, {a: 3, b: 1}, {a: 2, b: 2}, {a: 5, b: 2}] for Obj<a: Num, b: Num, Void> def <(other Obj<a: Num, b: Num, Void>) Bool as @b <(other @b) ok sort(<)`,
		`Arr<Obj<a: Num, b: Num, Void>>`,
		`[{a: 3, b: 1}, {a: 7, b: 2}, {a: 2, b: 2}, {a: 5, b: 2}]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[{a: 7}, {a: 3}, {a: 2}, {a: 5}] sortBy(@a, <)`,
		`Arr<Obj<a: Num, Void>>`,
		`[{a: 2}, {a: 3}, {a: 5}, {a: 7}]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[{a: 7, b: 2}, {a: 3, b: 1}, {a: 2, b: 2}, {a: 5, b: 2}] sortBy(@b, <)`,
		`Arr<Obj<a: Num, b: Num, Void>>`,
		`[{a: 3, b: 1}, {a: 7, b: 2}, {a: 2, b: 2}, {a: 5, b: 2}]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[] some`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[false] some`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[false, false] some`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[false, true, false] some`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[1, 3, 5, 2, 4, 7] takeWhile(if %2 ==1)`,
		`Arr<Num>`,
		`[1, 3, 5]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[1, 3, 5, 2, 4, 7] takeWhile(if %2 ==0)`,
		`Arr<Num>`,
		`[]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`[{a: 1}, {a: 2}, {b: 3}, {a: 4}] takeWhile(is {a: _}) each(@a)`,
		`Arr<Num>`,
		`[1, 2]`,
		nil,
		t,
	)
}
