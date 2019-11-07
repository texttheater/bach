package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestFilters(t *testing.T) {
	TestProgram(
		`["a", 1, "b", 2, "c", 3] each is Num with %2 >0 elis Str all`,
		&types.ArrType{types.Union(types.NumType{}, types.StrType{})},
		values.NewArrValue([]values.Value{
			values.StrValue("a"),
			values.NumValue(1),
			values.StrValue("b"),
			values.StrValue("c"),
			values.NumValue(3),
		}),
		nil,
		t,
	)
	TestProgram(
		`[{n: 1}, {n: 2}, {n: 3}] each is {n: n} with n %2 >0 all`,
		&types.ArrType{types.NewObjType(map[string]types.Type{
			"n": types.NumType{},
		})},
		values.NewArrValue([]values.Value{
			values.ObjValue(map[string]values.Value{
				"n": values.NumValue(1),
			}),
			values.ObjValue(map[string]values.Value{
				"n": values.NumValue(3),
			}),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3, 4, 5, 6] each if %2 ==0 then *2 else drop ok all`,
		&types.ArrType{types.NumType{}},
		values.NewArrValue([]values.Value{
			values.NumValue(4),
			values.NumValue(8),
			values.NumValue(12),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3, 4, 5, 6] each if %2 ==0 then drop else id ok all`,
		&types.ArrType{types.NumType{}},
		values.NewArrValue([]values.Value{
			values.NumValue(1),
			values.NumValue(3),
			values.NumValue(5),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3] each if %2 ==0 then drop else id ok +1 all`,
		&types.ArrType{types.NumType{}},
		values.NewArrValue([]values.Value{
			values.NumValue(2),
			values.NumValue(4),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3] each drop all`,
		&types.ArrType{types.VoidType{}},
		values.NewArrValue([]values.Value{}),
		nil,
		t,
	)
	TestProgram(
		`[{n: 1}, {n: 2}, {n: 3}] each is {n: n} then n ok all`,
		&types.ArrType{types.NumType{}},
		values.NewArrValue([]values.Value{
			values.NumValue(1),
			values.NumValue(2),
			values.NumValue(3),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3] each is Num all`,
		&types.ArrType{types.NumType{}},
		values.NewArrValue([]values.Value{
			values.NumValue(1),
			values.NumValue(2),
			values.NumValue(3),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3] each if ==1 then "a" elif ==2 then "b" else "c" ok all`,
		&types.ArrType{types.StrType{}},
		values.NewArrValue([]values.Value{
			values.StrValue("a"),
			values.StrValue("b"),
			values.StrValue("c"),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3] each if ==1 then "a" elif ==2 then "b" else "c" all`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.Syntax),
		),
		t,
	)
	TestProgram(
		`[1, 2, 3] each is Num ok +1 all`,
		&types.ArrType{types.NumType{}},
		values.NewArrValue([]values.Value{
			values.NumValue(2),
			values.NumValue(3),
			values.NumValue(4),
		}),
		nil,
		t,
	)
}
