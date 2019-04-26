package tests

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func UnionTypeTestCases() []TestCase {
	return []TestCase{
		{`for Num def f Num|Str as if ==1 then 1 else "abc" ok ok 1 f type`, types.StrType{}, values.StrValue("Num|Str"), nil},
		{`for Any def f Num|Str as 1 ok f type`, types.StrType{}, values.StrValue("Num|Str"), nil},
		{`for Any def f Void|Num as 1 ok f type`, types.StrType{}, values.StrValue("Num"), nil},
		{`for Any def f Num|Any as 1 ok f type`, types.StrType{}, values.StrValue("Any"), nil},
	}
}
