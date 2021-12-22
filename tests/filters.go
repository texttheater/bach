package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestFilters(t *testing.T) {
	TestProgram(
		`["a", 1, "b", 2, "c", 3] each is Num with %2 >0 elis Str all`,
		&types.Arr{types.NewUnion(types.Num{}, types.Str{})},
		states.NewArrValue([]states.Value{
			states.StrValue("a"),
			states.NumValue(1),
			states.StrValue("b"),
			states.StrValue("c"),
			states.NumValue(3),
		}),
		nil,
		t,
	)
	TestProgram(
		`[{n: 1}, {n: 2}, {n: 3}] each is {n: n} with n %2 >0 all`,
		&types.Arr{types.Obj{
			Props: map[string]types.Type{
				"n": types.Num{},
			},
			Rest: types.Void{},
		}},
		states.NewArrValue([]states.Value{
			states.ObjValueFromMap(map[string]states.Value{
				"n": states.NumValue(1),
			}),
			states.ObjValueFromMap(map[string]states.Value{
				"n": states.NumValue(3),
			}),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3, 4, 5, 6] each if %2 ==0 then *2 else drop ok all`,
		&types.Arr{types.Num{}},
		states.NewArrValue([]states.Value{
			states.NumValue(4),
			states.NumValue(8),
			states.NumValue(12),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3, 4, 5, 6] each if %2 ==0 then drop else id ok all`,
		&types.Arr{types.Num{}},
		states.NewArrValue([]states.Value{
			states.NumValue(1),
			states.NumValue(3),
			states.NumValue(5),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3] each if %2 ==0 then drop else id ok +1 all`,
		&types.Arr{types.Num{}},
		states.NewArrValue([]states.Value{
			states.NumValue(2),
			states.NumValue(4),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3] each drop all`,
		&types.Arr{types.Void{}},
		states.NewArrValue([]states.Value{}),
		nil,
		t,
	)
	TestProgram(
		`[{n: 1}, {n: 2}, {n: 3}] each is {n: n} then n ok all`,
		&types.Arr{types.Num{}},
		states.NewArrValue([]states.Value{
			states.NumValue(1),
			states.NumValue(2),
			states.NumValue(3),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3] each is Num all`,
		&types.Arr{types.Num{}},
		states.NewArrValue([]states.Value{
			states.NumValue(1),
			states.NumValue(2),
			states.NumValue(3),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3] each if ==1 then "a" elif ==2 then "b" else "c" ok all`,
		&types.Arr{types.Str{}},
		states.NewArrValue([]states.Value{
			states.StrValue("a"),
			states.StrValue("b"),
			states.StrValue("c"),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3] each if ==1 then "a" elif ==2 then "b" else "c" all`,
		nil,
		nil,
		errors.SyntaxError(
			errors.Code(errors.Syntax),
		),
		t,
	)
	TestProgram(
		`[1, 2, 3] each is Num ok +1 all`,
		&types.Arr{types.Num{}},
		states.NewArrValue([]states.Value{
			states.NumValue(2),
			states.NumValue(3),
			states.NumValue(4),
		}),
		nil,
		t,
	)
}
