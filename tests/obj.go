package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
)

func TestObjects(t *testing.T) {
	TestProgramStr(
		`{} get"a"`,
		`Void`,
		``,
		errors.TypeError(
			errors.Code(errors.VoidProgram),
		),
		t,
	)
	TestProgramStr(
		`{a: 1} get"a"`,
		`Num`,
		`1`,
		nil,
		t,
	)
	TestProgramStr(
		`{a: 1, b: "hey"} get"a"`,
		`Num|Str`,
		`1`,
		nil,
		t,
	)
	TestProgramStr(
		`{a: 1, b: "hey", c: false} get"a"`,
		`Num|Str|Bool`,
		`1`,
		nil,
		t,
	)
	TestProgramStr(
		`{1: "a"} get(1)`,
		`Str`,
		`"a"`,
		nil,
		t,
	)
	TestProgramStr(
		`{1.5: "a"} get(1.5)`,
		`Str`,
		`"a"`,
		nil,
		t,
	)
	TestProgramStr(
		`{a: 1, b: 2} props sort`,
		`Arr<Str>`,
		`["a", "b"]`,
		nil,
		t,
	)
	TestProgramStr(
		`{} props`,
		`Arr<Str>`,
		`[]`,
		nil,
		t,
	)
	TestProgramStr(
		`{a: 1, b: 2} values sort`,
		`Arr<Num>`,
		`[1, 2]`,
		nil,
		t,
	)
	TestProgramStr(
		`{} values`,
		`Tup<>`,
		`[]`,
		nil,
		t,
	)
	TestProgramStr(
		`{a: 1, b: 2} values sort`,
		`Arr<Num>`,
		`[1, 2]`,
		nil,
		t,
	)
	TestProgramStr(
		`{} values`,
		`Tup<>`,
		`[]`,
		nil,
		t,
	)
	TestProgramStr(
		`{a: 1, b: 2} items sortBy(@0)`,
		`Arr<Tup<Str, Num>>`,
		`[["a", 1], ["b", 2]]`,
		nil,
		t,
	)
	TestProgramStr(
		`{} items`,
		`Arr<Tup<Str, Void>>`,
		`[]`,
		nil,
		t,
	)
}
