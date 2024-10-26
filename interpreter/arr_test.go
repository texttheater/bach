package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestArrays(t *testing.T) {
	interpreter.TestProgram(
		`1 each(*2)`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.NoSuchFunction),
			errors.InputType(types.Num{}),
			errors.Name("each"),
			errors.NumParams(1),
		),
		t,
	)
	interpreter.TestProgram(
		`[1;2]`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.RestRequiresArrType),
			errors.WantType(types.AnyArr),
			errors.GotType(types.Num{}),
		),
		t,
	)
	interpreter.TestProgram(
		`[1 if ==2 then true else fatal ok] out`,
		nil,
		nil,
		errors.ValueError(
			errors.Code(errors.UnexpectedValue),
			errors.GotValue(states.NumValue(1)),
		),
		t,
	)
	interpreter.TestProgram(
		`[true, 1 if ==2 then true else fatal ok] out`,
		nil,
		nil,
		errors.ValueError(
			errors.Code(errors.UnexpectedValue),
			errors.GotValue(states.NumValue(1)),
		),
		t,
	)
}
