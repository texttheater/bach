package main_test

import (
	//"fmt"
	"math"
	//"os"
	"reflect"
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interp"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestType(t *testing.T) {
	cases := []struct {
		program string
		want types.Type
		errorKind string
	}{
		{`[1, 2, 3]`, &types.ArrType{&types.NumType{}}, ""},
	}
	for _, c := range cases {
		//fmt.Fprintf(os.Stderr, "%s\n", c.program)
		got, err := interp.TypecheckString(c.program)
		if c.errorKind != "" {
			if !errors.Is(c.errorKind, err) {
				t.Errorf("program: %v, want %v error, got %v error: %v", c.program, c.errorKind, errors.Kind(err), err)
			}
			continue
		}
		if err != nil {
			t.Errorf("program: %v, want %v, got %v error: %v", c.program, c.want, errors.Kind(err), err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("program: %v, want: %v, got: %v", c.program, c.want, got)
		}
	}
}

func TestValue(t *testing.T) {
	cases := []struct {
		program   string
		want      values.Value
		errorKind string
	}{
		// syntax errors
		{"&", nil, "syntax"},
		// type errors
		{"-1 *2", nil, "type"},
		{"3 <2 +1", nil, "type"},
		{"+", nil, "type"},
		{"hurz", nil, "type"},
		// literals
		{"1", &values.NumValue{1}, ""},
		{"1 2", &values.NumValue{2}, ""},
		{"1 2 3.5", &values.NumValue{3.5}, ""},
		// math
		{"1 +1", &values.NumValue{2}, ""},
		{"1 +2 *3", &values.NumValue{9}, ""},
		{"1 +(2 *3)", &values.NumValue{7}, ""},
		{"1 /0", &values.NumValue{math.Inf(1)}, ""},
		{"0 -1 *2", &values.NumValue{-2}, ""},
		{"15 %7", &values.NumValue{1}, ""},
		{"2 >3", &values.BoolValue{false}, ""},
		{"2 <3", &values.BoolValue{true}, ""},
		{"3 >2", &values.BoolValue{true}, ""},
		{"3 <2", &values.BoolValue{false}, ""},
		{"1 +1 ==2", &values.BoolValue{true}, ""},
		{"1 +1 >=2", &values.BoolValue{true}, ""},
		{"1 +1 <=2", &values.BoolValue{true}, ""},
		// logic
		{"true", &values.BoolValue{true}, ""},
		{"false", &values.BoolValue{false}, ""},
		{"true and(true)", &values.BoolValue{true}, ""},
		{"true and(false)", &values.BoolValue{false}, ""},
		{"false and(false)", &values.BoolValue{false}, ""},
		{"false and(true)", &values.BoolValue{false}, ""},
		{"true or(true)", &values.BoolValue{true}, ""},
		{"true or(false)", &values.BoolValue{true}, ""},
		{"false or(false)", &values.BoolValue{false}, ""},
		{"false or(true)", &values.BoolValue{true}, ""},
		{"false not", &values.BoolValue{true}, ""},
		{"true not", &values.BoolValue{false}, ""},
		{"true ==true", &values.BoolValue{true}, ""},
		{"true ==false", &values.BoolValue{false}, ""},
		{"false ==false", &values.BoolValue{true}, ""},
		{"false ==true", &values.BoolValue{false}, ""},
		{"1 +1 ==2 and(2 +2 ==5 not)", &values.BoolValue{true}, ""},
		// null
		{"1 null", &values.NullValue{}, ""},
		// assignment
		{"1 +1 =a 3 *2 +a", &values.NumValue{8}, ""},
		{"1 +1 ==2 =p 1 +1 ==1 =q p ==q not", &values.BoolValue{true}, ""},
		// strings
		{`"abc"`, &values.StrValue{"abc"}, ""},
		{`"\"\\abc\""`, &values.StrValue{`"\abc"`}, ""},
		{`1 "abc"`, &values.StrValue{"abc"}, ""},
		// arrays
		{`[]`, &values.ArrValue{[]values.Value{}}, ""},
		{`[1]`, &values.ArrValue{[]values.Value{&values.NumValue{1}}}, ""},
		{`[1, 2, 3]`, &values.ArrValue{[]values.Value{&values.NumValue{1}, &values.NumValue{2}, &values.NumValue{3}}}, ""},
		{`[1, "a"]`, &values.ArrValue{[]values.Value{&values.NumValue{1}, &values.StrValue{"a"}}}, ""},
		{`[[1, 2], ["a", "b"]]`, &values.ArrValue{[]values.Value{&values.ArrValue{[]values.Value{&values.NumValue{1}, &values.NumValue{2}}}, &values.ArrValue{[]values.Value{&values.StrValue{"a"}, &values.StrValue{"b"}}}}}, ""},
		// function definitions
		{`for Num def plusOne Num as +1 ok 1 plusOne`, &values.NumValue{2}, ""},
		{`for Num def plusOne Num as +1 ok 1 plusOne plusOne`, &values.NumValue{3}, ""},
		{`for Num def apply(for Num f Num) Num as f ok 1 apply(+1)`, &values.NumValue{2}, ""},
		{`for Num def connectSelf(for Num f(for Any g Num) Num) Num as =x f(x) ok 1 connectSelf(+)`, &values.NumValue{2}, ""},
		{`for Num def connectSelf(for Num f(for Any g Num) Num) Num as =x f(x) ok 1 connectSelf(+) 3 connectSelf(*)`, &values.NumValue{9}, ""},
		{`for Num def connectSelf(for Num f(g Num) Num) Num as =x f(x) ok 1 connectSelf(+)`, &values.NumValue{2}, ""},
		// bad function calls - TODO need finer-grained error values
		{`for Num def f Num as =x x ok for Str def f Str as =x x ok`, &values.NullValue{}, ""},
		{`for Num def f Num as =x x ok for Str def f Str as =x x ok f(2)`, nil, "type"}, // no such function (input type Any, name f, 1 parameters)
		{`for Num def f Num as =x x ok for Str def f Str as =x x ok 2 f`, &values.NumValue{2}, ""},
		{`for Num def f Num as =x x ok for Str def f Str as =x x ok f`, nil, "type"}, // no such function (input type Any, name f, 0 parameters)
		{`for Any def f(x Num) Num as x ok`, &values.NullValue{}, ""},
		{`for Any def f(x Num) Num as x ok f(1)`, &values.NumValue{1}, ""},
		{`for Any def f(x Num) Num as x ok for Any def f(x Str) Str as x ok`, &values.NullValue{}, ""},
		{`for Any def f(x Num) Num as x ok for Any def f(x Str) Str as x ok 1 f`, nil, "type"}, // no such function (input type Num, name f, 0 parameters
		{`for Any def f(x Num) Num as x ok for Any def f(x Str) Str as x ok f(1)`, nil, "type"}, // argument #0 has wrong output type (expected Str, got Num)
		{`for Any def f(for Num g Num) Num as 1 g ok`, &values.NullValue{}, ""},
		{`for Any def f(for Num g Num) Num as 1 g ok f(g)`, nil, "type"}, // no such function (input type Num, name g, 0 parameters)
		{`for Any def f(for Num g Num) Num as 1 g ok f(1)`, &values.NumValue{1}, ""},
		{`for Any def f(for Num g Num) Num as 1 g ok f(+1)`, &values.NumValue{2}, ""},
		{`for Any def f(for Num g Num) Num as 1 g ok f(+2)`, &values.NumValue{3}, ""},
		{`for Any def f(for Num g Num) Num as 1 g ok f(*10)`, &values.NumValue{10}, ""},
		{`for Any def f(for Num g(x Num) Num) Num as 1 g ok`, nil, "type"}, // no such function (input type Num, name g, 0 parameters
		{`for Any def f(for Num g(x Num) Num) Num as 1 g(2) ok`, &values.NullValue{}, ""},
		{`for Any def f(for Num g(x Num) Num) Num as 1 g(2) ok f(+)`, &values.NumValue{3}, ""},
		{`for Any def f(for Num g(x Num) Num) Num as 1 g(2) ok f(*)`, &values.NumValue{2}, ""},
		{`for Any def f(for Num g(x Num) Num) Num as 1 g(2) ok f(/)`, &values.NumValue{0.5}, ""},
		{`for Any def f(for Num g(x Num) Num) Num as 1 g(2) ok f(+1)`, nil, "type"}, // no such function (input type Num, name +, 2 parameters)
		{`for Any def f(for Num g(x Num) Num) Num as 1 g(2) ok for Any def g(x Str) Str as x ok`, &values.NullValue{}, ""},
		{`for Any def f(for Num g(x Num) Num) Num as 1 g(2) ok for Any def g(x Str) Str as x ok f(g)`, nil, "type"}, // argument #0 does not match parameters (expected x Num, got x Str)
		{`for Any def f(for Num g(x Num) Num) Num as 1 g(2) ok for Any def g(for Str x Num) Num as "abc" x ok f(g)`, nil, "type"}, // argument #0 does not match parameter (expected x Num, got for Str x Num)
		// conditionals
		{`if true then 2 else 3 ok`, &values.NumValue{2}, ""},
		{`for Num def heart Bool as if <3 then true else false ok ok 2 heart`, &values.BoolValue{true}, ""},
		{`for Num def heart Bool as if <3 then true else false ok ok 4 heart`, &values.BoolValue{false}, ""},
		{`for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 0 -1 expand`, &values.NumValue{-2}, ""},
		{`for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 1 expand`, &values.NumValue{2}, ""},
		{`for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 0 expand`, &values.NumValue{0}, ""},
		// recursion
		{`for Num def fac Num as if ==0 then 1 else =n -1 fac *n ok ok 3 fac`, &values.NumValue{6}, ""},
		// overloading
		{`for Bool def f Num as 1 ok for Num def f Num as 2 ok true f`, &values.NumValue{1}, ""},
		{`for Bool def f Num as 1 ok for Num def f Num as 2 ok 1 f`, &values.NumValue{2}, ""},
		// closures
		{`1 =a for Any def f Num as a ok f 2 =a f`, &values.NumValue{1}, ""},
		// sequences
		{`for Seq<Num> def f Seq<Num> as =x x ok [1, 2, 3] f`, &values.ArrValue{[]values.Value{&values.NumValue{1}, &values.NumValue{2}, &values.NumValue{3}}}, ""},
	}
	for _, c := range cases {
		//fmt.Fprintf(os.Stderr, "%s\n", c.program)
		got, err := interp.InterpretString(c.program)
		if c.errorKind != "" {
			if !errors.Is(c.errorKind, err) {
				t.Errorf("program: %v, want %v error, got %v error: %v", c.program, c.errorKind, errors.Kind(err), err)
			}
			continue
		}
		if err != nil {
			t.Errorf("program: %v, want %v, got %v error: %v", c.program, c.want, errors.Kind(err), err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("program: %v, want: %v, got: %v", c.program, c.want, got)
		}
	}
}
