package tests

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func StringTestCases() []TestCase {
	return []TestCase{
		{`"abc"`, types.StrType{}, values.StrValue("abc"), nil},
		{`"\"\\abc\""`, types.StrType{}, values.StrValue(`"\abc"`), nil},
		{`1 "abc"`, types.StrType{}, values.StrValue("abc"), nil},
	}
}
