package main_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/interp"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestInterp(t *testing.T) {
	cases := []struct {
		program   string
		wantType types.Type
		wantValue      values.Value
		wantErr error
	}{
		// syntax errors
		{"&", nil, nil, errors.E(errors.Kind(errors.Syntax))},
		// type errors
		{"-1 *2", nil, nil, errors.E(errors.Kind(errors.NoSuchFunction), errors.InputType(&types.NullType{}), errors.Name("-"), errors.NumParams(1))},
		{"3 <2 +1", nil, nil, errors.E(errors.Kind(errors.NoSuchFunction), errors.InputType(&types.BoolType{}), errors.Name("+"), errors.NumParams(1))},
		{"+", nil, nil, errors.E(errors.Kind(errors.NoSuchFunction), errors.InputType(&types.NullType{}), errors.Name("+"), errors.NumParams(0))},
		{"hurz", nil, nil, errors.E(errors.Kind(errors.NoSuchFunction))},
		// literals
		{"1", &types.NumType{}, &values.NumValue{1}, nil},
		{"1 2", &types.NumType{}, &values.NumValue{2}, nil},
		{"1 2 3.5", &types.NumType{}, &values.NumValue{3.5}, nil},
		// math
		{"1 +1", &types.NumType{}, &values.NumValue{2}, nil},
		{"1 +2 *3", &types.NumType{}, &values.NumValue{9}, nil},
		{"1 +(2 *3)", &types.NumType{}, &values.NumValue{7}, nil},
		{"1 /0", &types.NumType{}, &values.NumValue{math.Inf(1)}, nil},
		{"0 -1 *2", &types.NumType{}, &values.NumValue{-2}, nil},
		{"15 %7", &types.NumType{}, &values.NumValue{1}, nil},
		{"2 >3", &types.BoolType{}, &values.BoolValue{false}, nil},
		{"2 <3", &types.BoolType{}, &values.BoolValue{true}, nil},
		{"3 >2", &types.BoolType{}, &values.BoolValue{true}, nil},
		{"3 <2", &types.BoolType{}, &values.BoolValue{false}, nil},
		{"1 +1 ==2", &types.BoolType{}, &values.BoolValue{true}, nil},
		{"1 +1 >=2", &types.BoolType{}, &values.BoolValue{true}, nil},
		{"1 +1 <=2", &types.BoolType{}, &values.BoolValue{true}, nil},
		// logic
		{"true", &types.BoolType{}, &values.BoolValue{true}, nil},
		{"false", &types.BoolType{}, &values.BoolValue{false}, nil},
		{"true and(true)", &types.BoolType{}, &values.BoolValue{true}, nil},
		{"true and(false)", &types.BoolType{}, &values.BoolValue{false}, nil},
		{"false and(false)", &types.BoolType{}, &values.BoolValue{false}, nil},
		{"false and(true)", &types.BoolType{}, &values.BoolValue{false}, nil},
		{"true or(true)", &types.BoolType{}, &values.BoolValue{true}, nil},
		{"true or(false)", &types.BoolType{}, &values.BoolValue{true}, nil},
		{"false or(false)", &types.BoolType{}, &values.BoolValue{false}, nil},
		{"false or(true)", &types.BoolType{}, &values.BoolValue{true}, nil},
		{"false not", &types.BoolType{}, &values.BoolValue{true}, nil},
		{"true not", &types.BoolType{}, &values.BoolValue{false}, nil},
		{"true ==true", &types.BoolType{}, &values.BoolValue{true}, nil},
		{"true ==false", &types.BoolType{}, &values.BoolValue{false}, nil},
		{"false ==false", &types.BoolType{}, &values.BoolValue{true}, nil},
		{"false ==true", &types.BoolType{}, &values.BoolValue{false}, nil},
		{"1 +1 ==2 and(2 +2 ==5 not)", &types.BoolType{}, &values.BoolValue{true}, nil},
		// null
		{"1 null", &types.NullType{}, &values.NullValue{}, nil},
		// assignment
		{"1 +1 =a 3 *2 +a", &types.NumType{}, &values.NumValue{8}, nil},
		{"1 +1 ==2 =p 1 +1 ==1 =q p ==q not", &types.BoolType{}, &values.BoolValue{true}, nil},
		// strings
		{`"abc"`, &types.StrType{}, &values.StrValue{"abc"}, nil},
		{`"\"\\abc\""`, &types.StrType{}, &values.StrValue{`"\abc"`}, nil},
		{`1 "abc"`, &types.StrType{}, &values.StrValue{"abc"}, nil},
		// arrays
		{`[]`, &types.ArrType{&types.AnyType{}}, &values.ArrValue{[]values.Value{}}, nil},
		{`[1]`, &types.ArrType{&types.NumType{}}, &values.ArrValue{[]values.Value{&values.NumValue{1}}}, nil},
		{`[1, 2, 3]`, &types.ArrType{&types.NumType{}}, &values.ArrValue{[]values.Value{&values.NumValue{1}, &values.NumValue{2}, &values.NumValue{3}}}, nil},
		{`[1, "a"]`, &types.ArrType{&types.DisjunctiveType{[]types.Type{&types.NumType{}, &types.StrType{}}}}, &values.ArrValue{[]values.Value{&values.NumValue{1}, &values.StrValue{"a"}}}, nil},
		{`[[1, 2], ["a", "b"]]`, &types.ArrType{&types.DisjunctiveType{[]types.Type{&types.ArrType{&types.NumType{}}, &types.ArrType{&types.StrType{}}}}}, &values.ArrValue{[]values.Value{&values.ArrValue{[]values.Value{&values.NumValue{1}, &values.NumValue{2}}}, &values.ArrValue{[]values.Value{&values.StrValue{"a"}, &values.StrValue{"b"}}}}}, nil},
		// function definitions
		{`for Num def plusOne Num as +1 ok 1 plusOne`, &types.NumType{}, &values.NumValue{2}, nil},
		{`for Num def plusOne Num as +1 ok 1 plusOne plusOne`, &types.NumType{}, &values.NumValue{3}, nil},
		{`for Num def apply(for Num f Num) Num as f ok 1 apply(+1)`, &types.NumType{}, &values.NumValue{2}, nil},
		{`for Num def connectSelf(for Num f(for Any g Num) Num) Num as =x f(x) ok 1 connectSelf(+)`, &types.NumType{}, &values.NumValue{2}, nil},
		{`for Num def connectSelf(for Num f(for Any g Num) Num) Num as =x f(x) ok 1 connectSelf(+) 3 connectSelf(*)`, &types.NumType{}, &values.NumValue{9}, nil},
		{`for Num def connectSelf(for Num f(g Num) Num) Num as =x f(x) ok 1 connectSelf(+)`, &types.NumType{}, &values.NumValue{2}, nil},
		// bad function calls
		{`for Num def f Num as =x x ok for Str def f Str as =x x ok`, &types.NullType{}, &values.NullValue{}, nil},
		{`for Num def f Num as =x x ok for Str def f Str as =x x ok f(2)`, nil, nil, errors.E(errors.Kind(errors.NoSuchFunction), errors.InputType(&types.NullType{}), errors.Name("f"), errors.NumParams(1))},
		{`for Num def f Num as =x x ok for Str def f Str as =x x ok 2 f`, &types.NumType{}, &values.NumValue{2}, nil},
		{`for Num def f Num as =x x ok for Str def f Str as =x x ok f`, nil, nil, errors.E(errors.Kind(errors.NoSuchFunction), errors.InputType(&types.NullType{}), errors.Name("f"), errors.NumParams(0))},
		{`for Any def f(x Num) Num as x ok`, &types.NullType{}, &values.NullValue{}, nil},
		{`for Any def f(x Num) Num as x ok f(1)`, &types.NumType{}, &values.NumValue{1}, nil},
		{`for Any def f(x Num) Num as x ok for Any def f(x Str) Str as x ok`, &types.NullType{}, &values.NullValue{}, nil},
		{`for Any def f(x Num) Num as x ok for Any def f(x Str) Str as x ok 1 f`, nil, nil, errors.E(errors.Kind(errors.NoSuchFunction), errors.InputType(&types.NumType{}), errors.Name("f"), errors.NumParams(0))},
		{`for Any def f(x Num) Num as x ok for Any def f(x Str) Str as x ok f(1)`, nil, nil, errors.E(errors.Kind(errors.ArgHasWrongOutputType), errors.ArgNum(0), errors.WantType(&types.StrType{}), errors.GotType(&types.NumType{}))},
		{`for Any def f(for Num g Num) Num as 1 g ok`, &types.NullType{}, &values.NullValue{}, nil},
		{`for Any def f(for Num g Num) Num as 1 g ok f(g)`, nil, nil, errors.E(errors.Kind(errors.NoSuchFunction), errors.InputType(&types.NumType{}), errors.Name("g"), errors.NumParams(0))},
		{`for Any def f(for Num g Num) Num as 1 g ok f(1)`, &types.NumType{}, &values.NumValue{1}, nil},
		{`for Any def f(for Num g Num) Num as 1 g ok f(+1)`, &types.NumType{}, &values.NumValue{2}, nil},
		{`for Any def f(for Num g Num) Num as 1 g ok f(+2)`, &types.NumType{}, &values.NumValue{3}, nil},
		{`for Any def f(for Num g Num) Num as 1 g ok f(*10)`, &types.NumType{}, &values.NumValue{10}, nil},
		{`for Any def f(for Num g(x Num) Num) Num as 1 g ok`, nil, nil, errors.E(errors.Kind(errors.NoSuchFunction), errors.InputType(&types.NumType{}), errors.Name("g"), errors.NumParams(0))},
		{`for Any def f(for Num g(x Num) Num) Num as 1 g(2) ok`, &types.NullType{}, &values.NullValue{}, nil},
		{`for Any def f(for Num g(x Num) Num) Num as 1 g(2) ok f(+)`, &types.NumType{}, &values.NumValue{3}, nil},
		{`for Any def f(for Num g(x Num) Num) Num as 1 g(2) ok f(*)`, &types.NumType{}, &values.NumValue{2}, nil},
		{`for Any def f(for Num g(x Num) Num) Num as 1 g(2) ok f(/)`, &types.NumType{}, &values.NumValue{0.5}, nil},
		{`for Any def f(for Num g(x Num) Num) Num as 1 g(2) ok f(+1)`, nil, nil, errors.E(errors.Kind(errors.NoSuchFunction), errors.InputType(&types.NumType{}), errors.Name("+"), errors.NumParams(2))},
		{`for Any def f(for Num g(x Num) Num) Num as 1 g(2) ok for Any def g(x Str) Str as x ok`, &types.NullType{}, &values.NullValue{}, nil},
		{`for Any def f(for Num g(x Num) Num) Num as 1 g(2) ok for Any def g(x Str) Str as x ok f(g)`, nil, nil, errors.E(errors.Kind(errors.ParamDoesNotMatch), errors.ParamNum(0), errors.WantParam(&functions.Parameter{InputType: &types.AnyType{}, Name: "x", Params: nil, OutputType: &types.NumType{}}), errors.GotParam(&functions.Parameter{InputType: &types.AnyType{}, Name: "x", Params: nil, OutputType: &types.StrType{}}))},
		{`for Any def f(for Num g(x Num) Num) Num as 1 g(2) ok for Any def g(for Str x Num) Num as "abc" x ok f(g)`, nil, nil, errors.E(errors.Kind(errors.ParamDoesNotMatch), errors.ParamNum(0), errors.WantParam(&functions.Parameter{InputType: &types.AnyType{}, Name: "x", Params: nil, OutputType: &types.NumType{}}), errors.GotParam(&functions.Parameter{InputType: &types.StrType{}, Name: "x", Params: nil, OutputType: &types.NumType{}}))},
		// conditionals
		{`if true then 2 else 3 ok`, &types.NumType{}, &values.NumValue{2}, nil},
		{`for Num def heart Bool as if <3 then true else false ok ok 2 heart`,  &types.BoolType{}, &values.BoolValue{true}, nil},
		{`for Num def heart Bool as if <3 then true else false ok ok 4 heart`, &types.BoolType{}, &values.BoolValue{false}, nil},
		{`for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 0 -1 expand`, &types.NumType{}, &values.NumValue{-2}, nil},
		{`for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 1 expand`, &types.NumType{}, &values.NumValue{2}, nil},
		{`for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 0 expand`, &types.NumType{}, &values.NumValue{0}, nil},
		// recursion
		{`for Num def fac Num as if ==0 then 1 else =n -1 fac *n ok ok 3 fac`, &types.NumType{}, &values.NumValue{6}, nil},
		// overloading
		{`for Bool def f Num as 1 ok for Num def f Num as 2 ok true f`, &types.NumType{}, &values.NumValue{1}, nil},
		{`for Bool def f Num as 1 ok for Num def f Num as 2 ok 1 f`, &types.NumType{}, &values.NumValue{2}, nil},
		// closures
		{`1 =a for Any def f Num as a ok f 2 =a f`, &types.NumType{}, &values.NumValue{1}, nil},
		// sequences
		{`for Seq<Num> def f Seq<Num> as =x x ok [1, 2, 3] f`, &types.SeqType{&types.NumType{}}, &values.ArrValue{[]values.Value{&values.NumValue{1}, &values.NumValue{2}, &values.NumValue{3}}}, nil},
	}
	for _, c := range cases {
		gotType, gotValue, gotErr := interp.InterpretString(c.program)
		if c.wantErr != nil {
			if gotErr == nil {
				t.Log("ERROR: Expected error but program succeeded.")
				t.Logf("Program:        %s", c.program)
				t.Logf("Expected error: %s", c.wantErr)
				t.Logf("Got type:       %s", gotType)
				t.Logf("Got value:      %s", gotValue)
				t.Fail()
			} else if !errors.Match(c.wantErr, gotErr) {
				t.Log("ERROR: Expected error does not match actual error.")
				t.Logf("Program:        %s", c.program)
				t.Logf("Expected error: %s", c.wantErr)
				t.Logf("Got error:      %s", gotErr)
				t.Fail()
			}
		} else {
			if gotErr != nil {
				t.Log("ERROR: Expected program to succeed but got error.")
				t.Logf("Program:        %s", c.program)
				t.Logf("Expected type:  %s", c.wantType)
				t.Logf("Expected value: %s", c.wantValue)
				t.Logf("Got error:      %s", gotErr)
				t.Fail()
			} else if !reflect.DeepEqual(c.wantType, gotType) {
				t.Log("ERROR: Program has unexpected output type.")
				t.Logf("Program:        %s", c.program)
				t.Logf("Expected type:  %s", c.wantType)
				t.Logf("Expected value: %s", c.wantValue)
				t.Logf("Got type:       %s", gotType)
				t.Logf("Got value:      %s", gotValue)
				t.Fail()
			} else if !reflect.DeepEqual(c.wantValue, gotValue) {
				t.Log("ERROR: Program has unexpected output value.")
				t.Logf("Program:        %s", c.program)
				t.Logf("Expected type:  %s", c.wantType)
				t.Logf("Expected value: %s", c.wantValue)
				t.Logf("Got type:       %s", gotType)
				t.Logf("Got value:      %s", gotValue)
				t.Fail()
			}
		}
	}
}
