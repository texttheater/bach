package tests

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func SimpleTypeTestCases() []TestCase {
	return []TestCase{
		{`null type`, types.StrType{}, values.StrValue("Null"), nil},
		{`true type`, types.StrType{}, values.StrValue("Bool"), nil},
		{`1 type`, types.StrType{}, values.StrValue("Num"), nil},
		{`"abc" type`, types.StrType{}, values.StrValue("Str"), nil},
	}
}
