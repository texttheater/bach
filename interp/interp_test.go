package interp_test

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
		{"1", &values.NumberValue{1}, ""},
		{"1 2", &values.NumberValue{2}, ""},
		{"1 2 3.5", &values.NumberValue{3.5}, ""},
		{"1 +1", &values.NumberValue{2}, ""},
		{"1 +2 *3", &values.NumberValue{9}, ""},
		{"1 +(2 *3)", &values.NumberValue{7}, ""},
		{"1 /0", &values.NumberValue{math.Inf(1)}, ""},
		{"-1 *2", nil, "type"},
		{"0 -1 *2", &values.NumberValue{-2}, ""},
		//{"15 %7", &values.NumberValue{1}, ""},
		{"2 >3", &values.BooleanValue{false}, ""},
		{"2 <3", &values.BooleanValue{true}, ""},
		{"3 >2", &values.BooleanValue{true}, ""},
		{"3 <2", &values.BooleanValue{false}, ""},
		{"3 <2 +1", nil, "type"},
		{"+", nil, "syntax"},
	}
	for _, c := range cases {
		got, err := interp.InterpretString(c.program)
		if c.errorKind != "" {
			if !errors.Is(c.errorKind, err) {
				t.Errorf("program: %q, expected %q error, got: %q", c.program, c.errorKind, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("program: %q, want: %q, got error: %q", c.program, c.want, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("program: %q, want: %q, got: %q", c.program, c.want, got)
		}
	}
}
