package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestArrays(t *testing.T) {
	TestProgram(
		`[]`,
		types.NewTup([]types.Type{}),
		states.NewArrValue([]states.Value{}),
		nil,
		t,
	)
	TestProgram(
		`[1]`,
		types.NewTup([]types.Type{types.Num{}}),
		states.NewArrValue([]states.Value{states.NumValue(1)}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3]`,
		types.NewTup([]types.Type{types.Num{}, types.Num{}, types.Num{}}),
		states.NewArrValue([]states.Value{states.NumValue(1),
			states.NumValue(2),
			states.NumValue(3)}),
		nil,
		t,
	)
	TestProgram(
		`[1, "a"]`,
		types.NewTup([]types.Type{types.Num{}, types.Str{}}),
		states.NewArrValue([]states.Value{states.NumValue(1),
			states.StrValue("a")}),
		nil,
		t,
	)
	TestProgram(
		`[[1, 2], ["a", "b"]]`,
		types.NewTup([]types.Type{types.NewTup([]types.Type{types.Num{}, types.Num{}}), types.NewTup([]types.Type{types.Str{}, types.Str{}})}),
		states.NewArrValue([]states.Value{states.NewArrValue([]states.Value{states.NumValue(1),
			states.NumValue(2)}),
			states.NewArrValue([]states.Value{states.StrValue("a"), states.StrValue("b")})}),
		nil,
		t,
	)
	TestProgramStr(
		`for Arr<Num> def f Arr<Num> as =x x ok [1, 2, 3] f`,
		`Arr<Num>`,
		`[1, 2, 3]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] each(*2)`,
		`Arr<Num>`,
		`[2, 4, 6]`,
		nil,
		t,
	)
	TestProgram(
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
	TestProgram(
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
	TestProgram(
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
	TestProgram(
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
	TestProgram(
		`[1 if ==2 then true else fatal ok] out`,
		nil,
		nil,
		errors.ValueError(
			errors.Code(errors.UnexpectedValue),
			errors.GotValue(states.NumValue(1)),
		),
		t,
	)
	TestProgram(
		`[true, 1 if ==2 then true else fatal ok] out`,
		nil,
		nil,
		errors.ValueError(
			errors.Code(errors.UnexpectedValue),
			errors.GotValue(states.NumValue(1)),
		),
		t,
	)
	TestProgramStr(
		`[] get(0)`,
		`Void`,
		``,
		errors.TypeError(
			errors.Code(errors.VoidProgram),
		),
		t,
	)
	TestProgramStr(
		`["a", "b", "c"] get(0)`,
		`Str`,
		`"a"`,
		nil,
		t,
	)
	TestProgramStr(
		`["a", "b", "c"] get(-1)`,
		``,
		``,
		errors.ValueError(
			errors.Code(errors.BadIndex),
		),
		t,
	)
	TestProgramStr(
		`[1, 2, 3] take(2)`,
		`Arr<Num>`,
		`[1, 2]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] take(1)`,
		`Arr<Num>`,
		`[1]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] take(0)`,
		`Arr<Num>`,
		`[]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] take(-1)`,
		`Arr<Num>`,
		`[]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] take(4)`,
		`Arr<Num>`,
		`[1, 2, 3]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] take(3)`,
		`Arr<Num>`,
		`[1, 2, 3]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] drop(2)`,
		`Arr<Num>`,
		`[3]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] drop(1)`,
		`Arr<Num>`,
		`[2, 3]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] drop(0)`,
		`Arr<Num>`,
		`[1, 2, 3]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] drop(-1)`,
		`Arr<Num>`,
		`[1, 2, 3]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] drop(4)`,
		`Arr<Num>`,
		`[]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] drop(3)`,
		`Arr<Num>`,
		`[]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] rev`,
		`Arr<Num>`,
		`[3, 2, 1]`,
		nil,
		t,
	)
	TestProgramStr(
		`[] rev`,
		`Tup<>`,
		`[]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] +[4, 5]`,
		`Arr<Num>`,
		`[1, 2, 3, 4, 5]`,
		nil,
		t,
	)
	TestProgramStr(
		`[] +[4, 5]`,
		`Arr<Num>`,
		`[4, 5]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] +[]`,
		`Arr<Num>`,
		`[1, 2, 3]`,
		nil,
		t,
	)
	TestProgramStr(
		`[] +[]`,
		`Tup<>`,
		`[]`,
		nil,
		t,
	)
	TestProgramStr(
		`range(0, 4)`,
		`Arr<Num>`,
		`[0, 1, 2, 3]`,
		nil,
		t,
	)
	TestProgramStr(
		`range(-1, 2)`,
		`Arr<Num>`,
		`[-1, 0, 1]`,
		nil,
		t,
	)
	TestProgramStr(
		`range(3, 2)`,
		`Arr<Num>`,
		`[]`,
		nil,
		t,
	)
	TestProgramStr(
		`[7, 3, 2, 5, 2] sort`,
		`Arr<Num>`,
		`[2, 2, 3, 5, 7]`,
		nil,
		t,
	)
	TestProgramStr(
		`"Zwölf Boxkämpfer jagen Victor quer über den großen Sylter Deich . Voilà !" fields sort`,
		`Arr<Str>`,
		`["!", ".", "Boxkämpfer", "Deich", "Sylter", "Victor", "Voilà", "Zwölf", "den", "großen", "jagen", "quer", "über"]`,
		nil,
		t,
	)
	TestProgramStr(
		`"Zwölf Boxkämpfer jagen Victor quer über den großen Sylter Deich . Voilà !" fields sort(>)`,
		`Arr<Str>`,
		`["über", "quer", "jagen", "großen", "den", "Zwölf", "Voilà", "Victor", "Sylter", "Deich", "Boxkämpfer", ".", "!"]`,
		nil,
		t,
	)
	TestProgramStr(
		`[7, 3, 2, 5] sort(>)`,
		`Arr<Num>`,
		`[7, 5, 3, 2]`,
		nil,
		t,
	)
	TestProgramStr(
		`[{a: 7}, {a: 3}, {a: 2}, {a: 5}] for Obj<a: Num, Void> def <(other Obj<a: Num, Void>) Bool as @a <(other @a) ok sort(<)`,
		`Arr<Obj<a: Num, Void>>`,
		`[{a: 2}, {a: 3}, {a: 5}, {a: 7}]`,
		nil,
		t,
	)
	TestProgramStr(
		`[{a: 7, b: 2}, {a: 3, b: 1}, {a: 2, b: 2}, {a: 5, b: 2}] for Obj<a: Num, b: Num, Void> def <(other Obj<a: Num, b: Num, Void>) Bool as @b <(other @b) ok sort(<)`,
		`Arr<Obj<a: Num, b: Num, Void>>`,
		`[{a: 3, b: 1}, {a: 7, b: 2}, {a: 2, b: 2}, {a: 5, b: 2}]`,
		nil,
		t,
	)
	TestProgramStr(
		`[{a: 7}, {a: 3}, {a: 2}, {a: 5}] sortBy(@a, <)`,
		`Arr<Obj<a: Num, Void>>`,
		`[{a: 2}, {a: 3}, {a: 5}, {a: 7}]`,
		nil,
		t,
	)
	TestProgramStr(
		`[{a: 7, b: 2}, {a: 3, b: 1}, {a: 2, b: 2}, {a: 5, b: 2}] sortBy(@b, <)`,
		`Arr<Obj<a: Num, b: Num, Void>>`,
		`[{a: 3, b: 1}, {a: 7, b: 2}, {a: 2, b: 2}, {a: 5, b: 2}]`,
		nil,
		t,
	)
	TestProgramStr(
		`[] join`,
		`Arr<<A>>`, // FIXME we want Arr<Void>
		`[]`,
		nil,
		t,
	)
	TestProgramStr(
		`[[]] join`,
		`Arr<Void>`,
		`[]`,
		nil,
		t,
	)
	TestProgramStr(
		`[[1]] join`,
		`Arr<Num>`,
		`[1]`,
		nil,
		t,
	)
	TestProgramStr(
		`[[1, 2]] join`,
		`Arr<Num>`,
		`[1, 2]`,
		nil,
		t,
	)
	TestProgramStr(
		`[[], []] join`,
		`Arr<Void>`,
		`[]`,
		nil,
		t,
	)
	TestProgramStr(
		`[[], [1]] join`,
		`Arr<Num>`,
		`[1]`,
		nil,
		t,
	)
	TestProgramStr(
		`[[1], [2, 3]] join`,
		`Arr<Num>`,
		`[1, 2, 3]`,
		nil,
		t,
	)
	TestProgramStr(
		`[] every`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	TestProgramStr(
		`[true] every`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	TestProgramStr(
		`[true, true] every`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	TestProgramStr(
		`[true, false, true] every`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	TestProgramStr(
		`[] some`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	TestProgramStr(
		`[false] some`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	TestProgramStr(
		`[false, false] some`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	TestProgramStr(
		`[false, true, false] some`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	TestProgram(`[1, 2, 3] fold(0, +)`,
		types.Num{},
		states.NumValue(6),
		nil,
		t,
	)
	TestProgram(`[2, 3, 4] fold(1, *)`,
		types.Num{},
		states.NumValue(24),
		nil,
		t,
	)
	TestProgramStr(
		`[1, 3, 5, 2, 4, 7] takeWhile(if %2 ==1)`,
		`Arr<Num>`,
		`[1, 3, 5]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 3, 5, 2, 4, 7] takeWhile(if %2 ==0)`,
		`Arr<Num>`,
		`[]`,
		nil,
		t,
	)
	TestProgramStr(
		`[{a: 1}, {a: 2}, {b: 3}, {a: 4}] takeWhile(is {a: _}) each(@a)`,
		`Arr<Num>`,
		`[1, 2]`,
		nil,
		t,
	)
	TestProgramStr(
		`[{a: 1}, {a: 2}, {b: 3}, {a: 4}] dropWhile(is {a: _})`,
		`Arr<Obj<b: Num, Void>|Obj<a: Num, Void>>`,
		`[{b: 3}, {a: 4}]`,
		nil,
		t,
	)
}
