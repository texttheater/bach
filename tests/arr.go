package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestArrays(t *testing.T) {
	TestProgram(
		`[]`,
		types.TupType([]types.Type{}),
		values.NewArrValue([]values.Value{}),
		nil,
		t,
	)
	TestProgram(
		`[1]`,
		types.TupType([]types.Type{types.NumType{}}),
		values.NewArrValue([]values.Value{values.NumValue(1)}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3]`,
		types.TupType([]types.Type{types.NumType{}, types.NumType{}, types.NumType{}}),
		values.NewArrValue([]values.Value{values.NumValue(1),
			values.NumValue(2),
			values.NumValue(3)}),
		nil,
		t,
	)
	TestProgram(
		`[1, "a"]`,
		types.TupType([]types.Type{types.NumType{}, types.StrType{}}),
		values.NewArrValue([]values.Value{values.NumValue(1),
			values.StrValue("a")}),
		nil,
		t,
	)
	TestProgram(
		`[[1, 2], ["a", "b"]]`,
		types.TupType([]types.Type{types.TupType([]types.Type{types.NumType{}, types.NumType{}}), types.TupType([]types.Type{types.StrType{}, types.StrType{}})}),
		values.NewArrValue([]values.Value{values.NewArrValue([]values.Value{values.NumValue(1),
			values.NumValue(2)}),
			values.NewArrValue([]values.Value{values.StrValue("a"), values.StrValue("b")})}),
		nil,
		t,
	)
	TestProgram(
		`for Arr<Num> def f Arr<Num> as =x x ok [1, 2, 3] f`,
		&types.ArrType{
			types.NumType{},
		},
		values.NewArrValue(
			[]values.Value{
				values.NumValue(1),
				values.NumValue(2),
				values.NumValue(3),
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
		values.NewArrValue(
			[]values.Value{
				values.NumValue(2),
				values.NumValue(4),
				values.NumValue(6),
			},
		),
		nil,
		t,
	)
	TestProgram(
		`1 each *2 all`,
		nil,
		nil,
		errors.E(
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
		values.NewArrValue(
			[]values.Value{
				values.NumValue(1),
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
		values.NewArrValue(
			[]values.Value{
				values.NumValue(1),
				values.NumValue(2),
				values.NumValue(3),
				values.NumValue(4),
			},
		),
		nil,
		t,
	)
}
