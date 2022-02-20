package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestFilters(t *testing.T) {
	TestProgram(
		`["a", 1, "b", 2, "c", 3] keep(is Num with %2 >0 elis Str)`,
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
		`[{n: 1}, {n: 2}, {n: 3}] keep(is {n: n} with n %2 >0)`,
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
		`[1, 2, 3, 4, 5, 6] keep(if %2 ==0) each(*2)`,
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
		`[1, 2, 3, 4, 5, 6] keep(if %2 ==0 not) each(id)`,
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
		`[1, 2, 3] keep(if %2 ==0 not) each(+1)`,
		&types.Arr{types.Num{}},
		states.NewArrValue([]states.Value{
			states.NumValue(2),
			states.NumValue(4),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3] keep(if false)`,
		&types.Arr{types.Num{}},
		states.NewArrValue([]states.Value{}),
		nil,
		t,
	)
	TestProgram(
		`[{n: 1}, 2, {n: 3}] keep(is {n: n}) each(@n)`,
		&types.Arr{types.Num{}},
		states.NewArrValue([]states.Value{
			states.NumValue(1),
			states.NumValue(3),
		}),
		nil,
		t,
	)
	TestProgram(
		`[1, 2, 3] each(if ==1 then "a" elif ==2 then "b" else "c" ok)`,
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
		`[1, 2, 3] each(if ==1 then "a" elif ==2 then "b" else "c")`,
		nil,
		nil,
		errors.SyntaxError(
			errors.Code(errors.Syntax),
		),
		t,
	)
	TestProgram(
		`[1, 2, 3] each(+1)`,
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
