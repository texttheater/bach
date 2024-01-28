package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestMatchingType(t *testing.T) {
	interpreter.TestProgram(
		`if true then 2 else "two" ok is Num then true else false ok`,
		types.Bool{},
		states.BoolValue(true),
		nil,
		t,
	)
	//interpreter.TestProgram(
	//	`if true then 2 else "two" ok is Num then true ok`,
	//	nil,
	//	nil,
	//	errors.TypeError(
	//		errors.Code(errors.NonExhaustiveMatch),
	//		errors.WantType(types.VoidType{}),
	//		errors.GotType(types.StrType{}),
	//	),
	//	t,
	//)
	interpreter.TestProgram(
		`if true then 2 else "two" ok is Num then true elis Str then false ok`,
		types.Bool{},
		states.BoolValue(true),
		nil,
		t,
	)
	interpreter.TestProgram(
		`if true then 2 else "two" ok is Num then true elis Str then false else false ok`,
		nil,
		nil,
		errors.TypeError(
			errors.Code(errors.UnreachableElseClause),
		),
		t,
	)
	interpreter.TestProgramStr(
		`[1, 2, 3] is Arr<Num> then each(+1) ok`,
		`Arr<Num>`,
		`[2, 3, 4]`,
		nil,
		t,
	)
	interpreter.TestProgramStr( // Intersective matching: pattern says Arr<Any> but Bach knows it got Arr<Num>
		`[1, 2, 3] is Arr<Any> then each(+1) ok`,
		`Arr<Num>`,
		`[2, 3, 4]`,
		nil,
		t,
	)
}
