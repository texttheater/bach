package tests_test

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/tests"
	"github.com/texttheater/bach/types"
)

func TestObjectTypes(t *testing.T) {
	tests.TestProgram(`{} type`,
		types.Str{},
		states.StrValue("Obj<Void>"),
		nil,
		t,
	)
	tests.TestProgram(`{a: null} type`,
		types.Str{},
		states.StrValue("Obj<a: Null, Void>"),
		nil,
		t,
	)
	tests.TestProgram(`{b: false, a: null} type`,
		types.Str{},
		states.StrValue("Obj<a: Null, b: Bool, Void>"),
		nil,
		t,
	)
	tests.TestProgram(`{c: 0, b: false, a: null} type`,
		types.Str{},
		states.StrValue("Obj<a: Null, b: Bool, c: Num, Void>"),
		nil,
		t,
	)
	tests.TestProgram(`{d: "", c: 0, b: false, a: null} type`,
		types.Str{},
		states.StrValue("Obj<a: Null, b: Bool, c: Num, d: Str, Void>"),
		nil,
		t,
	)
	tests.TestProgram(`{e: [], d: "", c: 0, b: false, a: null} type`,
		types.Str{},
		states.StrValue("Obj<a: Null, b: Bool, c: Num, d: Str, e: Tup<>, Void>"),
		nil,
		t,
	)
	tests.TestProgram(`{f: {}, e: [], d: "", c: 0, b: false, a: null} type`,
		types.Str{},
		states.StrValue("Obj<a: Null, b: Bool, c: Num, d: Str, e: Tup<>, f: Obj<Void>, Void>"),
		nil,
		t,
	)
}
