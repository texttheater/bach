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
	}
	for _, c := range cases {
		got, err := interp.InterpretString(c.program)
		if c.errorKind != "" {
			if !errors.Is(c.errorKind, err) {
				t.Errorf("program: %q, want %v error, got %v error: %q", c.program, c.errorKind, errors.Kind(err), err)
			}
			continue
		}
		if err != nil {
			t.Errorf("program: %q, want %q, got %v error: %q", c.program, c.want, errors.Kind(err), err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("program: %q, want: %q, got: %q", c.program, c.want, got)
		}
	}
}
