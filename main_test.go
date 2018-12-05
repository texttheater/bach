package main_test

import (
	//"fmt"
	"math"
	//"os"
	"reflect"
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interp"
	"github.com/texttheater/bach/values"
)

func TestInterp(t *testing.T) {
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
		{"1", &values.NumberValue{1}, ""},
		{"1 2", &values.NumberValue{2}, ""},
		{"1 2 3.5", &values.NumberValue{3.5}, ""},
		// math
		{"1 +1", &values.NumberValue{2}, ""},
		{"1 +2 *3", &values.NumberValue{9}, ""},
		{"1 +(2 *3)", &values.NumberValue{7}, ""},
		{"1 /0", &values.NumberValue{math.Inf(1)}, ""},
		{"0 -1 *2", &values.NumberValue{-2}, ""},
		{"15 %7", &values.NumberValue{1}, ""},
		{"2 >3", &values.BooleanValue{false}, ""},
		{"2 <3", &values.BooleanValue{true}, ""},
		{"3 >2", &values.BooleanValue{true}, ""},
		{"3 <2", &values.BooleanValue{false}, ""},
		{"1 +1 ==2", &values.BooleanValue{true}, ""},
		{"1 +1 >=2", &values.BooleanValue{true}, ""},
		{"1 +1 <=2", &values.BooleanValue{true}, ""},
		// logic
		{"true", &values.BooleanValue{true}, ""},
		{"false", &values.BooleanValue{false}, ""},
		{"true and(true)", &values.BooleanValue{true}, ""},
		{"true and(false)", &values.BooleanValue{false}, ""},
		{"false and(false)", &values.BooleanValue{false}, ""},
		{"false and(true)", &values.BooleanValue{false}, ""},
		{"true or(true)", &values.BooleanValue{true}, ""},
		{"true or(false)", &values.BooleanValue{true}, ""},
		{"false or(false)", &values.BooleanValue{false}, ""},
		{"false or(true)", &values.BooleanValue{true}, ""},
		{"false not", &values.BooleanValue{true}, ""},
		{"true not", &values.BooleanValue{false}, ""},
		{"true ==true", &values.BooleanValue{true}, ""},
		{"true ==false", &values.BooleanValue{false}, ""},
		{"false ==false", &values.BooleanValue{true}, ""},
		{"false ==true", &values.BooleanValue{false}, ""},
		{"1 +1 ==2 and(2 +2 ==5 not)", &values.BooleanValue{true}, ""},
		// null
		{"1 null", &values.NullValue{}, ""},
		// assignment
		{"1 +1 =a 3 *2 +a", &values.NumberValue{8}, ""},
		{"1 +1 ==2 =p 1 +1 ==1 =q p ==q not", &values.BooleanValue{true}, ""},
		// strings
		{`"abc"`, &values.StringValue{"abc"}, ""},
		{`"\"\\abc\""`, &values.StringValue{`"\abc"`}, ""},
		{`1 "abc"`, &values.StringValue{"abc"}, ""},
		// arrays
		{`[]`, &values.ArrayValue{[]values.Value{}}, ""},
		{`[1]`, &values.ArrayValue{[]values.Value{&values.NumberValue{1}}}, ""},
		{`[1, 2, 3]`, &values.ArrayValue{[]values.Value{&values.NumberValue{1}, &values.NumberValue{2}, &values.NumberValue{3}}}, ""},
		{`[1, "a"]`, &values.ArrayValue{[]values.Value{&values.NumberValue{1}, &values.StringValue{"a"}}}, ""},
		{`[[1, 2], ["a", "b"]]`, &values.ArrayValue{[]values.Value{&values.ArrayValue{[]values.Value{&values.NumberValue{1}, &values.NumberValue{2}}}, &values.ArrayValue{[]values.Value{&values.StringValue{"a"}, &values.StringValue{"b"}}}}}, ""},
		// function definitions
		{`for Num def plusOne Num as +1 ok 1 plusOne`, &values.NumberValue{2}, ""},
		{`for Num def plusOne Num as +1 ok 1 plusOne plusOne`, &values.NumberValue{3}, ""},
		{`for Num def apply(for Num f Num) Num as f ok 1 apply(+1)`, &values.NumberValue{2}, ""},
		{`for Num def connectSelf(for Num f(for Any g Num) Num) Num as =x f(x) ok 1 connectSelf(+)`, &values.NumberValue{2}, ""},
		{`for Num def connectSelf(for Num f(for Any g Num) Num) Num as =x f(x) ok 1 connectSelf(+) 3 connectSelf(*)`, &values.NumberValue{9}, ""},
		{`for Num def connectSelf(for Num f(g Num) Num) Num as =x f(x) ok 1 connectSelf(+)`, &values.NumberValue{2}, ""},
		// conditionals
		{`if true then 2 else 3 ok`, &values.NumberValue{2}, ""},
		{`for Num def heart Bool as if <3 then true else false ok ok 2 heart`, &values.BooleanValue{true}, ""},
		{`for Num def heart Bool as if <3 then true else false ok ok 4 heart`, &values.BooleanValue{false}, ""},
		{`for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 0 -1 expand`, &values.NumberValue{-2}, ""},
		{`for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 1 expand`, &values.NumberValue{2}, ""},
		{`for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 0 expand`, &values.NumberValue{0}, ""},
		// recursion
		{`for Num def fac Num as if ==0 then 1 else =n -1 fac *n ok ok 3 fac`, &values.NumberValue{6}, ""},
		// overloading
		{`for Bool def f Num as 1 ok for Num def f Num as 2 ok true f`, &values.NumberValue{1}, ""},
		{`for Bool def f Num as 1 ok for Num def f Num as 2 ok 1 f`, &values.NumberValue{2}, ""},
		// closures
		{`1 =a for Any def f Num as a ok f 2 =a f`, &values.NumberValue{1}, ""},
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
