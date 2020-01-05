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
		types.NumType{},
		states.NumValue(2),
		nil,
		t,
	)
	TestProgram(
		`for Num def plusOne Num as +1 ok 1 plusOne plusOne`,
		types.NumType{},
		states.NumValue(3),
		nil,
		t,
	)
	TestProgram(
		`for Num def apply(for Num f Num) Num as f ok 1 apply(+1)`,
		types.NumType{},
		states.NumValue(2),
		nil,
		t,
	)
	TestProgram(
		`for Num def connectSelf(for Num f(for Any Num) Num) Num as =x f(x) ok 1 connectSelf(+)`,
		types.NumType{},
		states.NumValue(2),
		nil,
		t,
	)
	TestProgram(
		`for Num def connectSelf(for Num f(for Any Num) Num) Num as =x f(x) ok 1 connectSelf(+) 3 connectSelf(*)`,
		types.NumType{},
		states.NumValue(9),
		nil,
		t,
	)
	TestProgram(
		`for Num def connectSelf(for Num f(Num) Num) Num as =x f(x) ok 1 connectSelf(+)`,
		types.NumType{},
		states.NumValue(2),
		nil,
		t,
	)
	TestProgram(
		`for Num def apply(for Num f Num) Num as f ok 2 =n apply(+n)`,
		types.NumType{},
		states.NumValue(4),
		nil,
		t,
	)
	TestProgram(
		`for Any def f(a Num, b Num) Tup<Num, Num> as [a, b] ok f(2, 3)`,
		types.TupType([]types.Type{types.NumType{}, types.NumType{}}),
		states.NewArrValue([]states.Value{states.NumValue(2), states.NumValue(3)}),
		nil,
		t,
	)
	TestProgram(
		`for <A> def apply(for <A> f <B>) <B> as f ok 1 apply(+1)`,
		types.NumType{},
		states.NumValue(2),
		nil,
		t,
	)
	TestProgram(
		`for <A>|Null def myMust <A> as is Null then fatal else id ok ok null myMust`,
		types.TypeVariable{
			Name: "A",
		},
		nil,
		errors.E(
			errors.Code(errors.UnexpectedValue),
			errors.GotValue(states.NullValue{}),
		),
		t,
	)
	TestProgram(
		`for <A>|Null def myMust <A> as is Null then fatal else id ok ok 1 myMust`,
		types.NumType{},
		states.NumValue(1),
		nil,
		t,
	)
	TestProgram(
		`for <A>|Null def myMust <A> as is <A> then id else fatal ok ok null myMust`,
		types.NullType{},
		states.NullValue{},
		nil,
		t,
	)
	TestProgram(
		`for Str def f Obj<> as {} ok "abc" findFirst(f)`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.ArgHasWrongOutputType),
			errors.WantType(types.TypeVariable{
				Name: "A",
				UpperBound: types.Union(
					types.NullType{},
					types.NewObjType(map[string]types.Type{
						"start": types.NumType{},
						"0":     types.StrType{},
					})),
			},
			),
			errors.GotType(types.AnyObjType),
			errors.ArgNum(1),
		),
		t,
	)
	TestProgram(
		`for <A Obj<a: Num>> def f <A> as id ok {} f`,
		types.NewObjType(map[string]types.Type{
			"a": types.NumType{},
		}),
		nil,
		errors.E(
			errors.Code(errors.NoSuchFunction),
			errors.InputType(types.AnyObjType),
			errors.Name("f"),
			errors.NumParams(0),
		),
		t,
	)
	TestProgram(
		`for <A Obj<a: Num>> def f <A> as id ok {a: 1} f`,
		types.NewObjType(map[string]types.Type{
			"a": types.NumType{},
		}),
		states.ObjValue(map[string]*states.Thunk{
			"a": states.ThunkFromValue(states.NumValue(1)),
		}),
		nil,
		t,
	)
	// generics with bounds
	TestProgram(
		`for Any def f(for Any g <A Arr<Any>>) <A> as g ok f([1, "a"])`,
		types.TupType([]types.Type{
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
	TestProgram(
		`for Any def f(for Any g <A Arr<Any>>) <A> as g ok f("a")`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.ArgHasWrongOutputType),
			errors.ArgNum(1),
			errors.WantType(types.TypeVariable{"A", types.AnyArrType}),
			errors.GotType(types.StrType{}),
		),
		t,
	)
}
