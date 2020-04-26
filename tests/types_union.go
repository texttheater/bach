package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestUnionTypes(t *testing.T) {
	TestProgram(`for Num def f Num|Str as if ==1 then 1 else "abc" ok ok 1 f type`,
		types.StrType{},
		states.StrValue("Num|Str"),
		nil,
		t,
	)
	TestProgram(`for Any def f Num|Str as 1 ok f type`,
		types.StrType{},
		states.StrValue("Num|Str"),
		nil,
		t,
	)
	TestProgram(`for Any def f Void|Num as 1 ok f type`,
		types.StrType{},
		states.StrValue("Num"),
		nil,
		t,
	)
	TestProgram(`for Any def f Num|Any as 1 ok f type`,
		types.StrType{},
		states.StrValue("Any"),
		nil,
		t,
	)
}
