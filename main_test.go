package main_test

import (
	"math"
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
		// assignment
		{"1 +1 =a 3 *2 +a", &values.NumberValue{8}, ""},
		{"1 +1 ==2 =p 1 +1 ==1 =q p ==q not", &values.BooleanValue{true}, ""},
		// strings
		{`"abc"`, &values.StringValue{"abc"}, ""},
		{`"\"\\abc\""`, &values.StringValue{`"\abc"`}, ""},
		{`1 "abc"`, &values.StringValue{"abc"}, ""},
	}
	for _, c := range cases {
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
