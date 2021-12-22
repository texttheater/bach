package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestObjectTypes(t *testing.T) {
	TestProgram(`{} type`,
		types.Str{},
		states.StrValue("Obj<Void>"),
		nil,
		t,
	)
	TestProgram(`{a: null} type`,
		types.Str{},
		states.StrValue("Obj<a: Null, Void>"),
		nil,
		t,
	)
	TestProgram(`{b: false, a: null} type`,
		types.Str{},
		states.StrValue("Obj<a: Null, b: Bool, Void>"),
		nil,
		t,
	)
	TestProgram(`{c: 0, b: false, a: null} type`,
		types.Str{},
		states.StrValue("Obj<a: Null, b: Bool, c: Num, Void>"),
		nil,
		t,
	)
	TestProgram(`{d: "", c: 0, b: false, a: null} type`,
		types.Str{},
		states.StrValue("Obj<a: Null, b: Bool, c: Num, d: Str, Void>"),
		nil,
		t,
	)
	TestProgram(`{e: [], d: "", c: 0, b: false, a: null} type`,
		types.Str{},
		states.StrValue("Obj<a: Null, b: Bool, c: Num, d: Str, e: Tup<>, Void>"),
		nil,
		t,
	)
	TestProgram(`{f: {}, e: [], d: "", c: 0, b: false, a: null} type`,
		types.Str{},
		states.StrValue("Obj<a: Null, b: Bool, c: Num, d: Str, e: Tup<>, f: Obj<Void>, Void>"),
		nil,
		t,
	)
}
