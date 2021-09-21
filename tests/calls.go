package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/parameters"
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
		errors.TypeError(
			errors.Code(errors.NoSuchFunction),
			errors.InputType(types.NullType{}),
			errors.Name("f"),
			errors.NumParams(1),
		),
		t,
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
		errors.TypeError(
			errors.Code(errors.NoSuchFunction),
			errors.InputType(types.NullType{}),
			errors.Name("f"),
			errors.NumParams(0),
		),
		t,
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
		errors.TypeError(
			errors.Code(errors.NoSuchFunction),
			errors.InputType(types.NumType{}),
			errors.Name("f"),
			errors.NumParams(0),
		),
		t,
	)
	TestProgram(
		`for Any def f(x Num) Num as x ok for Any def f(x Str) Str as x ok f(1)`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.ArgHasWrongOutputType),
			errors.ArgNum(1),
			errors.WantType(types.StrType{}),
			errors.GotType(types.NumType{}),
		),
		t,
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
		errors.TypeError(
			errors.Code(errors.NoSuchFunction),
			errors.InputType(types.NumType{}),
			errors.Name("g"),
			errors.NumParams(0),
		),
		t,
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
		errors.TypeError(
			errors.Code(errors.NoSuchFunction),
			errors.InputType(types.NumType{}),
			errors.Name("g"),
			errors.NumParams(0),
		),
		t,
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
		errors.TypeError(
			errors.Code(errors.NoSuchFunction),
			errors.InputType(types.NumType{}),
			errors.Name("+"),
			errors.NumParams(2),
		),
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
		errors.TypeError(
			errors.Code(errors.ParamDoesNotMatch),
			errors.ParamNum(1),
			errors.WantParam(&parameters.Parameter{
				InputType:  types.AnyType{},
				Params:     nil,
				OutputType: types.NumType{},
			}),
			errors.GotParam(&parameters.Parameter{
				InputType:  types.AnyType{},
				Params:     nil,
				OutputType: types.StrType{},
			}),
		),
		t,
	)
	TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g(2) ok for Any def g(for Str x Num) Num as "abc" x ok f(g)`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.ParamDoesNotMatch),
			errors.ParamNum(1),
			errors.WantParam(&parameters.Parameter{
				InputType:  types.AnyType{},
				Params:     nil,
				OutputType: types.NumType{},
			}),
			errors.GotParam(&parameters.Parameter{
				InputType:  types.StrType{},
				Params:     nil,
				OutputType: types.NumType{},
			}),
		),
		t,
	)
}
