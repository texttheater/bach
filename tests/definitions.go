package tests

import (
	"testing"

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
	TestProgram(`for Num def apply(f for Num Num) Num as f ok 1 apply(+1)`,
		types.NumType{},
		values.NumValue(2),
		nil,
		t,
	)
	TestProgram(`for Num def connectSelf(f for Num (for Any Num) Num) Num as =x f(x) ok 1 connectSelf(+)`,
		types.NumType{},
		values.NumValue(2),
		nil,
		t,
	)
	TestProgram(`for Num def connectSelf(f for Num (for Any Num) Num) Num as =x f(x) ok 1 connectSelf(+) 3 connectSelf(*)`,
		types.NumType{},
		values.NumValue(9),
		nil,
		t,
	)
	TestProgram(`for Num def connectSelf(f for Num (Num) Num) Num as =x f(x) ok 1 connectSelf(+)`,
		types.NumType{},
		values.NumValue(2),
		nil,
		t,
	)
	TestProgram(`for <A> def apply(f for <A> <B>) <B> as f ok 1 apply(+1)`,
		types.NumType{},
		values.NumValue(2),
		nil,
		t,
	)
}
