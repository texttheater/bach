package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestFilters(t *testing.T) {
	TestProgram(
		`["a", 1, "b", 2, "c", 3] each is Num with %2 >0 elis Str all arr`,
		&types.ArrType{types.Union(types.NumType{}, types.StrType{})},
		values.ArrValue([]values.Value{
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
		`[{n: 1}, {n: 2}, {n: 3}] each is {n: n} with n %2 >0 all arr`,
		&types.ArrType{types.NewObjType(map[string]types.Type{
			"n": types.NumType{},
		})},
		values.ArrValue([]values.Value{
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
		`[{n: 1}, {n: 2}, {n: 3}] each is {n: n} then n ok all arr`,
		&types.ArrType{types.NumType{}},
		values.ArrValue([]values.Value{
			values.NumValue(1),
			values.NumValue(2),
			values.NumValue(3),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3] each is Num all arr`,
		&types.ArrType{types.NumType{}},
		values.ArrValue([]values.Value{
			values.NumValue(1),
			values.NumValue(2),
			values.NumValue(3),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3] each if ==1 then "a" elif ==2 then "b" else "c" ok all arr`,
		&types.ArrType{types.StrType{}},
		values.ArrValue([]values.Value{
			values.StrValue("a"),
			values.StrValue("b"),
			values.StrValue("c"),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3] each if ==1 then "a" elif ==2 then "b" else "c" all arr`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.Syntax),
		),
		t,
	)
	TestProgram(
		`[1, 2, 3] each is Num ok +1 all arr`,
		&types.ArrType{types.NumType{}},
		values.ArrValue([]values.Value{
			values.NumValue(2),
			values.NumValue(3),
			values.NumValue(4),
		}),
		nil,
		t,
	)
}
