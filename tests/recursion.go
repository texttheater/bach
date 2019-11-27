package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestRecursion(t *testing.T) {
	// naive recursive factorial
	TestProgram(
		`for Num def fac Num as if ==0 then 1 else =n -1 fac *n ok ok 3 fac`,
		types.NumType{},
		states.NumValue(6),
		nil,
		t,
	)
	// slightly different formulation where the recursive is in the
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
	TestProgram(
		`for Num def fac(acc Num) Num as =n if ==0 then acc else acc *n =acc n -1 fac(acc) ok ok 3 fac(1)`,
		types.NumType{},
		states.NumValue(6),
		nil,
		t,
	)
	// same, but with large input
	// this does not exhaust the goroutine stack
	// FIXME but it does seem to have a memory leak
	// (on the heap, then?)
	// it should be running in constant space
	// commented out because it's long-running
	//TestProgram(
	//	`for Num def fac(acc Num) Num as =n if ==0 then acc else acc *n =acc n -1 fac(acc) ok ok 10000000 fac(1)`,
	//	types.NumType{},
	//	states.NumValue(6),
	//	nil,
	//	t,
	//)
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
