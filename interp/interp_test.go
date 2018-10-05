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
		program string
		want    values.Value
	}{
		{"1", &values.NumberValue{1}},
		{"1 2", &values.NumberValue{2}},
		{"1 2 3.5", &values.NumberValue{3.5}},
		// TODO test for errors
		{"1 +1", &values.NumberValue{2}},
		{"1 +2 *3", &values.NumberValue{9}},
		{"1 +(2 *3)", &values.NumberValue{7}},
		{"1 /0", &values.NumberValue{math.Inf(1)}},
		{"0 -1 *2", &values.NumberValue{-2}},
	}
	for _, c := range cases {
		got, err := interp.InterpretString(c.program)
		if err != nil {
			errors.Explain(err, c.program)
			t.Errorf("program: %q, want: %q, got: error", c.program, c.want)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("program: %q, want: %q, got: %q", c.program, c.want, got)
		}
	}
}
