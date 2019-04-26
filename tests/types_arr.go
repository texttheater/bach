package tests

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func ArrayTypeTestCases() []TestCase {
	return []TestCase{
		{`[] type`, types.StrType{}, values.StrValue("Tup<>"), nil},
		{`["dog", "cat"] type`, types.StrType{}, values.StrValue("Tup<Str, Str>"), nil},
		{`["dog", 1] type`, types.StrType{}, values.StrValue("Tup<Str, Num>"), nil},
		{`["dog", 1, {}] type`, types.StrType{}, values.StrValue("Tup<Str, Num, Obj<>>"), nil},
		{`["dog", 1, {}, 2] type`, types.StrType{}, values.StrValue("Tup<Str, Num, Obj<>, Num>"), nil},
	}
}
