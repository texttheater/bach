package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestCalls(t *testing.T) {
	interpreter.TestProgram(
		`for Num def f Num as =x x ok for Str def f Str as =x x ok`,
		types.Null{},
		states.NullValue{},
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Num def f Num as =x x ok for Str def f Str as =x x ok f(2)`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.NoSuchFuncer),
			errors.InputType(types.Null{}),
			errors.Name("f"),
			errors.NumParams(1),
		),
		t,
	)
	interpreter.TestProgram(
		`for Num def f Num as =x x ok for Str def f Str as =x x ok 2 f`,
		types.Num{},
		states.NumValue(2),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Num def f Num as =x x ok for Str def f Str as =x x ok f`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.NoSuchFuncer),
			errors.InputType(types.Null{}),
			errors.Name("f"),
			errors.NumParams(0),
		),
		t,
	)
	interpreter.TestProgram(
		`for Any def f(x Num) Num as x ok`,
		types.Null{},
		states.NullValue{},
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Any def f(x Num) Num as x ok f(1)`,
		types.Num{},
		states.NumValue(1),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Any def f(x Num) Num as x ok for Any def f(x Str) Str as x ok`,
		types.Null{},
		states.NullValue{},
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Any def f(x Num) Num as x ok for Any def f(x Str) Str as x ok 1 f`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.NoSuchFuncer),
			errors.InputType(types.Num{}),
			errors.Name("f"),
			errors.NumParams(0),
		),
		t,
	)
	interpreter.TestProgram(
		`for Any def f(x Num) Num as x ok for Any def f(x Str) Str as x ok f(1)`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.ArgHasWrongOutputType),
			errors.ArgNum(1),
			errors.WantType(types.Str{}),
			errors.GotType(types.Num{}),
		),
		t,
	)
	interpreter.TestProgram(
		`for Any def f(for Num g Num) Num as 1 g ok`,
		types.Null{},
		states.NullValue{},
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Any def f(for Num g Num) Num as 1 g ok f(g)`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.NoSuchFuncer),
			errors.InputType(types.Num{}),
			errors.Name("g"),
			errors.NumParams(0),
		),
		t,
	)
	interpreter.TestProgram(
		`for Any def f(for Num g Num) Num as 1 g ok f(1)`,
		types.Num{},
		states.NumValue(1),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Any def f(for Num g Num) Num as 1 g ok f(+1)`,
		types.Num{},
		states.NumValue(2),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Any def f(for Num g Num) Num as 1 g ok f(+2)`,
		types.Num{},
		states.NumValue(3),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Any def f(for Num g Num) Num as 1 g ok f(*10)`,
		types.Num{},
		states.NumValue(10),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g ok`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.NoSuchFuncer),
			errors.InputType(types.Num{}),
			errors.Name("g"),
			errors.NumParams(0),
		),
		t,
	)
	interpreter.TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g(2) ok`,
		types.Null{},
		states.NullValue{},
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g(2) ok f(+)`,
		types.Num{},
		states.NumValue(3),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g(2) ok f(*)`,
		types.Num{},
		states.NumValue(2),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g(2) ok f(/)`,
		types.Num{},
		states.NumValue(0.5),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g(2) ok f(+1)`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.NoSuchFuncer),
			errors.InputType(types.Num{}),
			errors.Name("+"),
			errors.NumParams(2),
		),
		t,
	)
	interpreter.TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g(2) ok for Any def g(x Str) Str as x ok`,
		types.Null{},
		states.NullValue{},
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g(2) ok for Any def g(x Str) Str as x ok f(g)`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.ParamDoesNotMatch),
			errors.ParamNum(1),
			errors.GotParam(params.SimpleParam("", "", types.Str{})),
			errors.WantParam(params.SimpleParam("", "", types.Num{})),
		),
		t,
	)
	interpreter.TestProgram(
		`for Any def f(for Num g(Num) Num) Num as 1 g(2) ok for Any def g(for Str x Num) Num as "abc" x ok f(g)`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.ParamDoesNotMatch),
			errors.ParamNum(1),
			errors.GotParam(&params.Param{
				InputType:  types.Str{},
				OutputType: types.Num{},
			}),
			errors.WantParam(params.SimpleParam("", "", types.Num{})),
		),
		t,
	)
	interpreter.TestProgramStr(
		`a[2]`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.NoSuchFuncer),
			errors.InputType(types.Null{}),
			errors.Name(`a`),
			errors.NumParams(1),
		),
		t,
	)
	interpreter.TestProgramStr(
		`a{b: 2}`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.NoSuchFuncer),
			errors.InputType(types.Null{}),
			errors.Name(`a`),
			errors.NumParams(1),
		),
		t,
	)
	interpreter.TestProgramStr(
		`for Num def applyWith2AsArg(for Num f(Num) <A>) <A> as f(2) ok 1 applyWith2AsArg(+)`,
		`Num`,
		`3`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`for Num def applyWith2AsArg(for Num f(Num) <A>) <A> as f(2) ok 1 applyWith2AsArg(*)`,
		`Num`,
		`2`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`for Num def applyWithABCAsArg(for Num f(Str) <A>) <A> as f("abc") ok 1 applyWithABCAsArg(+)`,
		``,
		``,
		errors.TypeError(
			errors.Code(errors.ParamDoesNotMatch),
			errors.ParamNum(1),
			errors.WantParam(params.SimpleParam("", "", types.Str{})),
			errors.GotParam(params.SimpleParam("", "", types.Num{})),
		),
		t,
	)
}
