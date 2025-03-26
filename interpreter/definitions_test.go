package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestDefinitions(t *testing.T) {
	interpreter.TestProgram(
		`for Num def plusOne Num as +1 ok 1 plusOne`,
		types.NumType{},
		states.NumValue(2),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Num def plusOne Num as +1 ok 1 plusOne plusOne`,
		types.NumType{},
		states.NumValue(3),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Num def apply(for Num f Num) Num as f ok 1 apply(+1)`,
		types.NumType{},
		states.NumValue(2),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Num def connectSelf(for Num f(for Any Num) Num) Num as =x f(x) ok 1 connectSelf(+)`,
		types.NumType{},
		states.NumValue(2),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Num def connectSelf(for Num f(for Any Num) Num) Num as =x f(x) ok 1 connectSelf(+) 3 connectSelf(*)`,
		types.NumType{},
		states.NumValue(9),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Num def connectSelf(for Num f(Num) Num) Num as =x f(x) ok 1 connectSelf(+)`,
		types.NumType{},
		states.NumValue(2),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Num def apply(for Num f Num) Num as f ok 2 =n apply(+n)`,
		types.NumType{},
		states.NumValue(4),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Any def f(a Num, b Num) Arr<Num, Num> as [a, b] ok f(2, 3)`,
		types.NewTup([]types.Type{types.NumType{}, types.NumType{}}),
		states.NewArrValue([]states.Value{states.NumValue(2), states.NumValue(3)}),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for <A> def apply(for <A> f <B>) <B> as f ok 1 apply(+1)`,
		types.NumType{},
		states.NumValue(2),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for <A>|Null def myMust <A> as is Null then fatal else id ok ok null myMust`,
		types.TypeVar{
			Name: "A",
		},
		nil,
		errors.ValueError(
			errors.Code(errors.UnexpectedValue),
			errors.GotValue(states.NullValue{}),
		),
		t,
	)
	interpreter.TestProgram(
		`for <A>|Null def myMust <A> as is Null then fatal else id ok ok 1 myMust`,
		types.NumType{},
		states.NumValue(1),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for <A>|Null def myMust <A> as is <A> then id else fatal ok ok null myMust`,
		types.NullType{},
		states.NullValue{},
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Str def f Obj<> as {} ok "abc" reFindAll(f)`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.ArgHasWrongOutputType),
			errors.WantType(types.TypeVar{
				Name: "A",
				Bound: types.NewUnionType(
					types.NullType{},
					types.ObjType{
						Props: map[string]types.Type{
							"start": types.NumType{},
							"0":     types.StrType{},
						},
						Rest: types.AnyType{},
					},
				),
			},
			),
			errors.GotType(types.AnyObjType),
			errors.ArgNum(1),
		),
		t,
	)
	interpreter.TestProgram(
		`for <A Obj<a: Num>> def f <A> as id ok {} f`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.NoSuchFunction),
			errors.InputType(types.VoidObjType),
			errors.Name("f"),
			errors.NumParams(0),
		),
		t,
	)
	interpreter.TestProgram(
		`for <A Obj<a: Num, Any>> def f <A> as id ok {a: 1} f`,
		types.ObjType{
			Props: map[string]types.Type{
				"a": types.NumType{},
			},
			Rest: types.VoidType{},
		},
		states.ObjValue(map[string]*states.Thunk{
			"a": states.ThunkFromValue(states.NumValue(1)),
		}),
		nil,
		t,
	)
	// generics with bounds
	interpreter.TestProgram(
		`for Any def f(for Any g <A Arr<Any...>>) <A> as g ok f([1, "a"])`,
		types.NewTup([]types.Type{
			types.NumType{},
			types.StrType{},
		}),
		states.NewArrValue([]states.Value{
			states.NumValue(1),
			states.StrValue("a"),
		}),
		nil,
		t,
	)
	interpreter.TestProgram(
		`for Any def f(for Any g <A Arr<Any...>>) <A> as g ok f("a")`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.ArgHasWrongOutputType),
			errors.ArgNum(1),
			errors.WantType(types.NewTypeVar("A", types.AnyArrType)),
			errors.GotType(types.StrType{}),
		),
		t,
	)
}
