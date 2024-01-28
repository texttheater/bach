package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
)

func TestObjects(t *testing.T) {
	interpreter.TestProgramStr(
		`{} get"a"`,
		`Void`,
		``,
		errors.TypeError(
			errors.Code(errors.VoidProgram),
		),
		t,
	)
	interpreter.TestProgramStr(
		`{a: 1} get"a"`,
		`Num`,
		`1`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`{a: 1, b: "hey"} get"a"`,
		`Num|Str`,
		`1`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`{a: 1, b: "hey", c: false} get"a"`,
		`Num|Str|Bool`,
		`1`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`{1: "a"} get(1)`,
		`Str`,
		`"a"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`{1.5: "a"} get(1.5)`,
		`Str`,
		`"a"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`{a: 1, b: 2} props sort`,
		`Arr<Str>`,
		`["a", "b"]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`{} props`,
		`Arr<Str>`,
		`[]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`{a: 1, b: 2} values sort`,
		`Arr<Num>`,
		`[1, 2]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`{} values`,
		`Tup<>`,
		`[]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`{a: 1, b: 2} values sort`,
		`Arr<Num>`,
		`[1, 2]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`{} values`,
		`Tup<>`,
		`[]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`{a: 1, b: 2} items sortByStr(@0)`,
		`Arr<Tup<Str, Num>>`,
		`[["a", 1], ["b", 2]]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`{} items`,
		`Arr<Tup<Str, Void>>`,
		`[]`,
		nil,
		t,
	)
}
