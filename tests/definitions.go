package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestDefinitions(t *testing.T) {
	TestProgram(
		`for Num def plusOne Num as +1 ok 1 plusOne`,
		types.Num{},
		states.NumValue(2),
		nil,
		t,
	)
	TestProgram(
		`for Num def plusOne Num as +1 ok 1 plusOne plusOne`,
		types.Num{},
		states.NumValue(3),
		nil,
		t,
	)
	TestProgram(
		`for Num def apply(for Num f Num) Num as f ok 1 apply(+1)`,
		types.Num{},
		states.NumValue(2),
		nil,
		t,
	)
	TestProgram(
		`for Num def connectSelf(for Num f(for Any Num) Num) Num as =x f(x) ok 1 connectSelf(+)`,
		types.Num{},
		states.NumValue(2),
		nil,
		t,
	)
	TestProgram(
		`for Num def connectSelf(for Num f(for Any Num) Num) Num as =x f(x) ok 1 connectSelf(+) 3 connectSelf(*)`,
		types.Num{},
		states.NumValue(9),
		nil,
		t,
	)
	TestProgram(
		`for Num def connectSelf(for Num f(Num) Num) Num as =x f(x) ok 1 connectSelf(+)`,
		types.Num{},
		states.NumValue(2),
		nil,
		t,
	)
	TestProgram(
		`for Num def apply(for Num f Num) Num as f ok 2 =n apply(+n)`,
		types.Num{},
		states.NumValue(4),
		nil,
		t,
	)
	TestProgram(
		`for Any def f(a Num, b Num) Tup<Num, Num> as [a, b] ok f(2, 3)`,
		types.NewTup([]types.Type{types.Num{}, types.Num{}}),
		states.NewArrValue([]states.Value{states.NumValue(2), states.NumValue(3)}),
		nil,
		t,
	)
	TestProgram(
		`for <A> def apply(for <A> f <B>) <B> as f ok 1 apply(+1)`,
		types.Num{},
		states.NumValue(2),
		nil,
		t,
	)
	TestProgram(
		`for <A>|Null def myMust <A> as is Null then fatal else id ok ok null myMust`,
		types.Var{
			Name: "A",
		},
		nil,
		errors.ValueError(
			errors.Code(errors.UnexpectedValue),
			errors.GotValue(states.NullValue{}),
		),
		t,
	)
	TestProgram(
		`for <A>|Null def myMust <A> as is Null then fatal else id ok ok 1 myMust`,
		types.Num{},
		states.NumValue(1),
		nil,
		t,
	)
	TestProgram(
		`for <A>|Null def myMust <A> as is <A> then id else fatal ok ok null myMust`,
		types.Null{},
		states.NullValue{},
		nil,
		t,
	)
	TestProgram(
		`for Str def f Obj<> as {} ok "abc" findAll(f)`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.ArgHasWrongOutputType),
			errors.WantType(types.Var{
				Name: "A",
				Bound: types.NewUnion(
					types.Null{},
					types.Obj{
						Props: map[string]types.Type{
							"start": types.Num{},
							"0":     types.Str{},
						},
						Rest: types.Any{},
					},
				),
			},
			),
			errors.GotType(types.AnyObj),
			errors.ArgNum(1),
		),
		t,
	)
	TestProgram(
		`for <A Obj<a: Num>> def f <A> as id ok {} f`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.NoSuchFunction),
			errors.InputType(types.VoidObj),
			errors.Name("f"),
			errors.NumParams(0),
		),
		t,
	)
	TestProgram(
		`for <A Obj<a: Num, Any>> def f <A> as id ok {a: 1} f`,
		types.Obj{
			Props: map[string]types.Type{
				"a": types.Num{},
			},
			Rest: types.Void{},
		},
		states.ObjValue(map[string]*states.Thunk{
			"a": states.ThunkFromValue(states.NumValue(1)),
		}),
		nil,
		t,
	)
	// generics with bounds
	TestProgram(
		`for Any def f(for Any g <A Arr<Any>>) <A> as g ok f([1, "a"])`,
		types.NewTup([]types.Type{
			types.Num{},
			types.Str{},
		}),
		states.NewArrValue([]states.Value{
			states.NumValue(1),
			states.StrValue("a"),
		}),
		nil,
		t,
	)
	TestProgram(
		`for Any def f(for Any g <A Arr<Any>>) <A> as g ok f("a")`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.ArgHasWrongOutputType),
			errors.ArgNum(1),
			errors.WantType(types.NewVar("A", types.AnyArr)),
			errors.GotType(types.Str{}),
		),
		t,
	)
}
