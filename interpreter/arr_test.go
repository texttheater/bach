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
