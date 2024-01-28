package tests_test

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/tests"
	"github.com/texttheater/bach/types"
)

func TestUnionTypes(t *testing.T) {
	tests.TestProgram(`for Num def f Num|Str as if ==1 then 1 else "abc" ok ok 1 f type`,
		types.Str{},
		states.StrValue("Num|Str"),
		nil,
		t,
	)
	tests.TestProgram(`for Any def f Num|Str as 1 ok f type`,
		types.Str{},
		states.StrValue("Num|Str"),
		nil,
		t,
	)
	tests.TestProgram(`for Any def f Void|Num as 1 ok f type`,
		types.Str{},
		states.StrValue("Num"),
		nil,
		t,
	)
	tests.TestProgram(`for Any def f Num|Any as 1 ok f type`,
		types.Str{},
		states.StrValue("Any"),
		nil,
		t,
	)
}
