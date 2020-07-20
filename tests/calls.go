package tests

import (
	"testing"

	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestCalls(t *testing.T) {
	TestProgram(
		`for Num def f Num as =x x ok for Str def f Str as =x x ok`,
		types.NullType{},
		states.NullValue{},
		nil,
		t,
	)
	TestProgram(
		`for Num def f Num as =x x ok for Str def f Str as =x x ok f(2)`,
		nil,
		nil,
		states.E(states.Code(states.NoSuchFunction), states.InputType(types.NullType{}), states.Name("f"), states.NumParams(1)), t,
	)
	TestProgram(
		`for Num def f Num as =x x ok for Str def f Str as =x x ok 2 f`,
		types.NumType{},
		states.NumValue(2),
		nil,
		t,
	)
	TestProgram(
		`for Num def f Num as =x x ok for Str def f Str as =x x ok f`,
		nil,
		nil,
		states.E(states.Code(states.NoSuchFunction), states.InputType(types.NullType{}), states.Name("f"), states.NumParams(0)), t,
	)
	TestProgram(
		`for Any def f(x Num) Num as x ok`,
		types.NullType{},
		states.NullValue{},
		nil,
		t,
	)
	TestProgram(
		`for Any def f(x Num) Num as x ok f(1)`,
		types.NumType{},
		states.NumValue(1),
		nil,
		t,
	)
	TestProgram(
		`for Any def f(x Num) Num as x ok for Any def f(x Str) Str as x ok`,
		types.NullType{},
		states.NullValue{},
		nil,
		t,
	)
	TestProgram(
		`for Any def f(x Num) Num as x ok for Any def f(x Str) Str as x ok 1 f`,
		nil,
		nil,
		states.E(states.Code(states.NoSuchFunction), states.InputType(types.NumType{}), states.Name("f"), states.NumParams(0)), t,
	)
	TestProgram(
		`for Any def f(x Num) Num as x ok for Any def f(x Str) Str as x ok f(1)`,
		nil,
		nil,
		states.E(states.Code(states.ArgHasWrongOutputType), states.ArgNum(1), states.WantType(types.StrType{}), states.GotType(types.NumType{})), t,
	)
	TestProgram(
		`for Any def f(for Num g Num) Num as 1 g ok`,
		types.NullType{},
		states.NullValue{},
		nil,
		t,
	)
	TestProgram(
		`for Any def f(for Num g Num) Num as 1 g ok f(g)`,
		nil,
		nil,
		states.E(states.Code(states.NoSuchFunction), states.InputType(types.NumType{}), states.Name("g"), states.NumParams(0)), t,
	)
	TestProgram(
		`for Any def f(for Num g Num) Num as 1 g ok f(1)`,
		types.NumType{},
		states.NumValue(1),
		nil,
		t,
	)
	TestProgram(
		`for Any def f(for Num g Num) Num as 1 g ok f(+1)`,
		types.NumType{},
		states.NumValue(2),
		nil,
		t,
	)
	TestProgram(
		`for Any def f(for Num g Num) Num as 1 g ok f(+2)`,
		types.NumType{},
		states.NumValue(3),
		nil,
		t,
	)
	TestProgram(
		`for Any def f(for Num g Num) Num as 1 g ok f(*10)`,
		types.NumType{},
		states.NumValue(10),
		nil,
		t,
	)
	TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g ok`,
		nil,
		nil,
		states.E(states.Code(states.NoSuchFunction), states.InputType(types.NumType{}), states.Name("g"), states.NumParams(0)), t,
	)
	TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g(2) ok`,
		types.NullType{},
		states.NullValue{},
		nil,
		t,
	)
	TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g(2) ok f(+)`,
		types.NumType{},
		states.NumValue(3),
		nil,
		t,
	)
	TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g(2) ok f(*)`,
		types.NumType{},
		states.NumValue(2),
		nil,
		t,
	)
	TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g(2) ok f(/)`,
		types.NumType{},
		states.NumValue(0.5),
		nil,
		t,
	)
	TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g(2) ok f(+1)`,
		nil,
		nil,
		states.E(
			states.Code(states.NoSuchFunction),
			states.InputType(types.NumType{}),
			states.Name("+"),
			states.NumParams(2)),

		t,
	)
	TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g(2) ok for Any def g(x Str) Str as x ok`,
		types.NullType{},
		states.NullValue{},
		nil,
		t,
	)
	TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g(2) ok for Any def g(x Str) Str as x ok f(g)`,
		nil,
		nil,
		states.E(
			states.Code(states.ParamDoesNotMatch),
			states.ParamNum(1),
			states.WantParam(&functions.Parameter{
				InputType:  types.AnyType{},
				Params:     nil,
				OutputType: types.NumType{},
			}),
			states.GotParam(&functions.Parameter{
				InputType:  types.AnyType{},
				Params:     nil,
				OutputType: types.StrType{},
			})),

		t,
	)
	TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g(2) ok for Any def g(for Str x Num) Num as "abc" x ok f(g)`,
		nil,
		nil,
		states.E(
			states.Code(states.ParamDoesNotMatch),
			states.ParamNum(1),
			states.WantParam(&functions.Parameter{
				InputType:  types.AnyType{},
				Params:     nil,
				OutputType: types.NumType{},
			}),
			states.GotParam(&functions.Parameter{
				InputType:  types.StrType{},
				Params:     nil,
				OutputType: types.NumType{},
			})),

		t,
	)
}
