package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestDefinitions(t *testing.T) {
	TestProgram(
		`for Num def plusOne Num as +1 ok 1 plusOne`,
		types.NumType{},
		values.NumValue(2),
		nil,
		t,
	)
	TestProgram(`for Num def plusOne Num as +1 ok 1 plusOne plusOne`,
		types.NumType{},
		values.NumValue(3),
		nil,
		t,
	)
	TestProgram(`for Num def apply(for Num f Num) Num as f ok 1 apply(+1)`,
		types.NumType{},
		values.NumValue(2),
		nil,
		t,
	)
	TestProgram(`for Num def connectSelf(for Num f(for Any Num) Num) Num as =x f(x) ok 1 connectSelf(+)`,
		types.NumType{},
		values.NumValue(2),
		nil,
		t,
	)
	TestProgram(`for Num def connectSelf(for Num f(for Any Num) Num) Num as =x f(x) ok 1 connectSelf(+) 3 connectSelf(*)`,
		types.NumType{},
		values.NumValue(9),
		nil,
		t,
	)
	TestProgram(`for Num def connectSelf(for Num f(Num) Num) Num as =x f(x) ok 1 connectSelf(+)`,
		types.NumType{},
		values.NumValue(2),
		nil,
		t,
	)
	TestProgram(`for Num def apply(for Num f Num) Num as f ok 2 =n apply(+n)`,
		types.NumType{},
		values.NumValue(4),
		nil,
		t,
	)
	TestProgram(`for Any def f(a Num, b Num) Tup<Num, Num> as [a, b] ok f(2, 3)`,
		types.TupType([]types.Type{types.NumType{}, types.NumType{}}),
		values.NewArrValue([]values.Value{values.NumValue(2), values.NumValue(3)}),
		nil,
		t,
	)
	TestProgram(`for <A> def apply(for <A> f <B>) <B> as f ok 1 apply(+1)`,
		types.NumType{},
		values.NumValue(2),
		nil,
		t,
	)
	TestProgram(`for <A>|Null def must <A> as is Null then reject else id ok ok null must`,
		types.TypeVariable{
			Name: "A",
		},
		nil,
		errors.E(
			errors.Code(errors.UnexpectedValue),
			errors.GotValue(values.NullValue{}),
		),
		t,
	)
	TestProgram(`for <A>|Null def must <A> as is Null then reject else id ok ok 1 must`,
		types.NumType{},
		values.NumValue(1),
		nil,
		t,
	)
	TestProgram(`for <A>|Null def must <A> as is <A> then id else reject ok ok null must`,
		types.NullType{},
		values.NullValue{},
		nil,
		t,
	)
}
