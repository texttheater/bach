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
		types.TupType([]types.Type{}),
		states.NewArrValue([]states.Value{}),
		nil,
		t,
	)
	TestProgram(
		`[1]`,
		types.TupType([]types.Type{types.NumType{}}),
		states.NewArrValue([]states.Value{states.NumValue(1)}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3]`,
		types.TupType([]types.Type{types.NumType{}, types.NumType{}, types.NumType{}}),
		states.NewArrValue([]states.Value{states.NumValue(1),
			states.NumValue(2),
			states.NumValue(3)}),
		nil,
		t,
	)
	TestProgram(
		`[1, "a"]`,
		types.TupType([]types.Type{types.NumType{}, types.StrType{}}),
		states.NewArrValue([]states.Value{states.NumValue(1),
			states.StrValue("a")}),
		nil,
		t,
	)
	TestProgram(
		`[[1, 2], ["a", "b"]]`,
		types.TupType([]types.Type{types.TupType([]types.Type{types.NumType{}, types.NumType{}}), types.TupType([]types.Type{types.StrType{}, types.StrType{}})}),
		states.NewArrValue([]states.Value{states.NewArrValue([]states.Value{states.NumValue(1),
			states.NumValue(2)}),
			states.NewArrValue([]states.Value{states.StrValue("a"), states.StrValue("b")})}),
		nil,
		t,
	)
	TestProgram(
		`for Arr<Num> def f Arr<Num> as =x x ok [1, 2, 3] f`,
		&types.ArrType{
			types.NumType{},
		},
		states.NewArrValue(
			[]states.Value{
				states.NumValue(1),
				states.NumValue(2),
				states.NumValue(3),
			},
		),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3] each *2 all`,
		&types.ArrType{
			types.NumType{},
		},
		states.NewArrValue(
			[]states.Value{
				states.NumValue(2),
				states.NumValue(4),
				states.NumValue(6),
			},
		),
		nil,
		t,
	)
	TestProgram(
		`1 each *2 all`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.MappingRequiresArrType),
		),
		t,
	)
	TestProgram(
		`[1;[]]`,
		types.TupType(
			[]types.Type{
				types.NumType{},
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
		types.TupType(
			[]types.Type{
				types.NumType{},
				types.NumType{},
				types.NumType{},
				types.NumType{},
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
			errors.WantType(types.AnyArrType),
			errors.GotType(types.NumType{}),
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
		`[1, 2, 3] slice(0, 1)`,
		`Arr<Num>`,
		`[1]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] slice(0, 2)`,
		`Arr<Num>`,
		`[1, 2]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] slice(0, 3)`,
		`Arr<Num>`,
		`[1, 2, 3]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] slice(0, 4)`,
		`Arr<Num>`,
		`[1, 2, 3]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] slice(1, 3)`,
		`Arr<Num>`,
		`[2, 3]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] slice(2, 3)`,
		`Arr<Num>`,
		`[3]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] slice(3, 3)`,
		`Arr<Num>`,
		`[]`,
		nil,
		t,
	)
	TestProgramStr(
		`[1, 2, 3] slice(4, 3)`,
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
}
