package tests

import (
	"testing"

	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestFilters(t *testing.T) {
	TestProgram(
		`["a", 1, "b", 2, "c", 3] eachis Num with %2 >0 elis Str all arr`,
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
		`[{n: 1}, {n: 2}, {n: 3}] eachis {n: n} with n %2 >0 all arr`,
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
		`[{n: 1}, {n: 2}, {n: 3}] eachis {n: n} then n all arr`,
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
		`[1, 2, 3] eachis Num all arr`,
		&types.ArrType{types.NumType{}},
		values.ArrValue([]values.Value{
			values.NumValue(1),
			values.NumValue(2),
			values.NumValue(3),
		}),
		nil,
		t,
	)
}
