package tests

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func AnyTypeTestCases() []TestCase {
	return []TestCase{
		{`for Any def f Any as null ok f type`, types.StrType{}, values.StrValue("Any"), nil},
	}
}
