package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestRecursion(t *testing.T) {
	// simplest tail recursion example
	TestProgram(
		`for Num def f Num as if ==0 then 0 else -1 f ok ok 10000 f`,
		types.NumType{},
		states.NumValue(0),
		nil,
		t,
	)
	// naive recursive factorial
	TestProgram(
		`for Num def fac Num as if ==0 then 1 else =n -1 fac *n ok ok 3 fac`,
		types.NumType{},
		states.NumValue(6),
		nil,
		t,
	)
	// slightly different formulation where the recursive call is in the
	// argument of *
	TestProgram(
		`for Num def fac Num as if ==0 then 1 else =n *(n -1 fac) ok ok 3 fac`,
		types.NumType{},
		states.NumValue(6),
		nil,
		t,
	)
	// unorthodox factorial where the input value is ignored and n is
	// passed as an argument instead
	TestProgram(
		`for Any def fac(n Num) Num as n if ==0 then 1 else fac(n -1) *n ok ok fac(3)`,
		types.NumType{},
		states.NumValue(6),
		nil,
		t,
	)
	// tail-recursive factorial
	// This does not exhaust the goroutine stack and runs in constant space.
	TestProgram(
		`for Num def fac(acc Num) Num as =n if ==0 then acc else acc *n =acc n -1 fac(acc) ok ok 3 fac(1)`,
		types.NumType{},
		states.NumValue(6),
		nil,
		t,
	)
	// fold for numbers
	TestProgram(`for Arr<Num> def fold(start Num, for Num op(Num) Num) Num as is [head;tail] then start op(head) =newStart tail fold(newStart, op) else start ok ok [1, 2, 3] fold(0, +)`,
		types.NumType{},
		states.NumValue(6),
		nil,
		t,
	)
	// generic fold
	TestProgram(`for Arr<<A>> def fold(start <B>, for <B> op(<A>) <B>) <B> as is [head;tail] then start op(head) =newStart tail fold(newStart, op) else start ok ok [2, 3, 4] fold(1, *)`,
		types.NumType{},
		states.NumValue(24),
		nil,
		t,
	)
}
