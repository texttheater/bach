package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
)

func TestConditionals(t *testing.T) {
	TestProgramStr(`if true then 2 else 3 ok`,
		`Num`,
		`2`,
		nil,
		t,
	)
	TestProgramStr(`for Num def heart Bool as if <3 then true else false ok ok 2 heart`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	TestProgramStr(`for Num def heart Bool as if <3 then true else false ok ok 4 heart`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	TestProgramStr(`for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 0 -1 expand`,
		`Num`,
		`-2`,
		nil,
		t,
	)
	TestProgramStr(`for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 1 expand`,
		`Num`,
		`2`,
		nil,
		t,
	)
	TestProgramStr(`for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 0 expand`,
		`Num`,
		`0`,
		nil,
		t,
	)
	// predicates
	TestProgramStr(
		`is Null`,
		``,
		``,
		errors.E(
			errors.Code(errors.UnreachableElseClause),
		),
		t,
	)
	TestProgramStr(
		`2 is Num with >3`,
		`Obj<yes: Num>|Obj<no: Num>`,
		`{no: 2}`,
		nil,
		t,
	)
	TestProgramStr(
		`4 is Num with >3`,
		`Obj<yes: Num>|Obj<no: Num>`,
		`{yes: 4}`,
		nil,
		t,
	)
	TestProgramStr(
		`2 if >3`,
		`Obj<yes: Num>|Obj<no: Num>`,
		`{no: 2}`,
		nil,
		t,
	)
	TestProgramStr(
		`4 if >3`,
		`Obj<yes: Num>|Obj<no: Num>`,
		`{yes: 4}`,
		nil,
		t,
	)
	TestProgramStr(
		`for Any def f Num|Str as 2 ok f is Num _`,
		`Obj<yes: Num>|Obj<no: Str>`,
		`{yes: 2}`,
		nil,
		t,
	)
}
