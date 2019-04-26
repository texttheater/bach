package tests

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func ObjectTypeTestCases() []TestCase {
	return []TestCase{
		{`{} type`, types.StrType{}, values.StrValue("Obj<>"), nil},
		{`{a: null} type`, types.StrType{}, values.StrValue("Obj<a: Null>"), nil},
		{`{b: false, a: null} type`, types.StrType{}, values.StrValue("Obj<a: Null, b: Bool>"), nil},
		{`{c: 0, b: false, a: null} type`, types.StrType{}, values.StrValue("Obj<a: Null, b: Bool, c: Num>"), nil},
		{`{d: "", c: 0, b: false, a: null} type`, types.StrType{}, values.StrValue("Obj<a: Null, b: Bool, c: Num, d: Str>"), nil},
		{`{e: [], d: "", c: 0, b: false, a: null} type`, types.StrType{}, values.StrValue("Obj<a: Null, b: Bool, c: Num, d: Str, e: Tup<>>"), nil},
		{`{f: {}, e: [], d: "", c: 0, b: false, a: null} type`, types.StrType{}, values.StrValue("Obj<a: Null, b: Bool, c: Num, d: Str, e: Tup<>, f: Obj<>>"), nil},
	}
}
