package tests

import (
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func CallTestCases() []TestCase {
	return []TestCase{
		{`for Num def f Num as =x x ok for Str def f Str as =x x ok`, types.NullType{}, &values.NullValue{}, nil},
		{`for Num def f Num as =x x ok for Str def f Str as =x x ok f(2)`, nil, nil, errors.E(errors.Code(errors.NoSuchFunction), errors.InputType(types.NullType{}), errors.Name("f"), errors.NumParams(1))},
		{`for Num def f Num as =x x ok for Str def f Str as =x x ok 2 f`, types.NumType{}, values.NumValue(2), nil},
		{`for Num def f Num as =x x ok for Str def f Str as =x x ok f`, nil, nil, errors.E(errors.Code(errors.NoSuchFunction), errors.InputType(types.NullType{}), errors.Name("f"), errors.NumParams(0))},
		{`for Any def f(x Num) Num as x ok`, types.NullType{}, &values.NullValue{}, nil},
		{`for Any def f(x Num) Num as x ok f(1)`, types.NumType{}, values.NumValue(1), nil},
		{`for Any def f(x Num) Num as x ok for Any def f(x Str) Str as x ok`, types.NullType{}, &values.NullValue{}, nil},
		{`for Any def f(x Num) Num as x ok for Any def f(x Str) Str as x ok 1 f`, nil, nil, errors.E(errors.Code(errors.NoSuchFunction), errors.InputType(types.NumType{}), errors.Name("f"), errors.NumParams(0))},
		{`for Any def f(x Num) Num as x ok for Any def f(x Str) Str as x ok f(1)`, nil, nil, errors.E(errors.Code(errors.ArgHasWrongOutputType), errors.ArgNum(0), errors.WantType(types.StrType{}), errors.GotType(types.NumType{}))},
		{`for Any def f(g for Num Num) Num as 1 g ok`, types.NullType{}, &values.NullValue{}, nil},
		{`for Any def f(g for Num Num) Num as 1 g ok f(g)`, nil, nil, errors.E(errors.Code(errors.NoSuchFunction), errors.InputType(types.NumType{}), errors.Name("g"), errors.NumParams(0))},
		{`for Any def f(g for Num Num) Num as 1 g ok f(1)`, types.NumType{}, values.NumValue(1), nil},
		{`for Any def f(g for Num Num) Num as 1 g ok f(+1)`, types.NumType{}, values.NumValue(2), nil},
		{`for Any def f(g for Num Num) Num as 1 g ok f(+2)`, types.NumType{}, values.NumValue(3), nil},
		{`for Any def f(g for Num Num) Num as 1 g ok f(*10)`, types.NumType{}, values.NumValue(10), nil},
		{`for Any def f(g for Num (Num) Num) Num as 1 g ok`, nil, nil, errors.E(errors.Code(errors.NoSuchFunction), errors.InputType(types.NumType{}), errors.Name("g"), errors.NumParams(0))},
		{`for Any def f(g for Num (Num) Num) Num as 1 g(2) ok`, types.NullType{}, &values.NullValue{}, nil},
		{`for Any def f(g for Num (Num) Num) Num as 1 g(2) ok f(+)`, types.NumType{}, values.NumValue(3), nil},
		{`for Any def f(g for Num (Num) Num) Num as 1 g(2) ok f(*)`, types.NumType{}, values.NumValue(2), nil},
		{`for Any def f(g for Num (Num) Num) Num as 1 g(2) ok f(/)`, types.NumType{}, values.NumValue(0.5), nil},
		{`for Any def f(g for Num (Num) Num) Num as 1 g(2) ok f(+1)`, nil, nil, errors.E(errors.Code(errors.NoSuchFunction), errors.InputType(types.NumType{}), errors.Name("+"), errors.NumParams(2))},
		{`for Any def f(g for Num (Num) Num) Num as 1 g(2) ok for Any def g(x Str) Str as x ok`, types.NullType{}, &values.NullValue{}, nil},
		{`for Any def f(g for Num (Num) Num) Num as 1 g(2) ok for Any def g(x Str) Str as x ok f(g)`, nil, nil, errors.E(errors.Code(errors.ParamDoesNotMatch), errors.ParamNum(0), errors.WantParam(&shapes.Parameter{InputType: types.AnyType{}, Params: nil, OutputType: types.NumType{}}), errors.GotParam(&shapes.Parameter{InputType: types.AnyType{}, Params: nil, OutputType: types.StrType{}}))},
		{`for Any def f(g for Num (Num) Num) Num as 1 g(2) ok for Any def g(x for Str Num) Num as "abc" x ok f(g)`, nil, nil, errors.E(errors.Code(errors.ParamDoesNotMatch), errors.ParamNum(0), errors.WantParam(&shapes.Parameter{InputType: types.AnyType{}, Params: nil, OutputType: types.NumType{}}), errors.GotParam(&shapes.Parameter{InputType: types.StrType{}, Params: nil, OutputType: types.NumType{}}))},
	}
}
