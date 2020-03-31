package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestObjectTypes(t *testing.T) {
	TestProgram(`{} type`,
		types.StrType{},
		states.StrValue("Obj<>"),
		nil,
		t,
	)
	TestProgram(`{a: null} type`,
		types.StrType{},
		states.StrValue("Obj<a: Null>"),
		nil,
		t,
	)
	TestProgram(`{b: false, a: null} type`,
		types.StrType{},
		states.StrValue("Obj<a: Null, b: Bool>"),
		nil,
		t,
	)
	TestProgram(`{c: 0, b: false, a: null} type`,
		types.StrType{},
		states.StrValue("Obj<a: Null, b: Bool, c: Num>"),
		nil,
		t,
	)
	TestProgram(`{d: "", c: 0, b: false, a: null} type`,
		types.StrType{},
		states.StrValue("Obj<a: Null, b: Bool, c: Num, d: Str>"),
		nil,
		t,
	)
	TestProgram(`{e: [], d: "", c: 0, b: false, a: null} type`,
		types.StrType{},
		states.StrValue("Obj<a: Null, b: Bool, c: Num, d: Str, e: Tup<>>"),
		nil,
		t,
	)
	TestProgram(`{f: {}, e: [], d: "", c: 0, b: false, a: null} type`,
		types.StrType{},
		states.StrValue("Obj<a: Null, b: Bool, c: Num, d: Str, e: Tup<>, f: Obj<>>"),
		nil,
		t,
	)
}
