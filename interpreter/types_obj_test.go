package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestObjectTypes(t *testing.T) {
	interpreter.TestProgram(`{} type`,
		types.Str{},
		states.StrValue("Obj<Void>"),
		nil,
		t,
	)
	interpreter.TestProgram(`{a: null} type`,
		types.Str{},
		states.StrValue("Obj<a: Null, Void>"),
		nil,
		t,
	)
	interpreter.TestProgram(`{b: false, a: null} type`,
		types.Str{},
		states.StrValue("Obj<a: Null, b: Bool, Void>"),
		nil,
		t,
	)
	interpreter.TestProgram(`{c: 0, b: false, a: null} type`,
		types.Str{},
		states.StrValue("Obj<a: Null, b: Bool, c: Num, Void>"),
		nil,
		t,
	)
	interpreter.TestProgram(`{d: "", c: 0, b: false, a: null} type`,
		types.Str{},
		states.StrValue("Obj<a: Null, b: Bool, c: Num, d: Str, Void>"),
		nil,
		t,
	)
	interpreter.TestProgram(`{e: [], d: "", c: 0, b: false, a: null} type`,
		types.Str{},
		states.StrValue("Obj<a: Null, b: Bool, c: Num, d: Str, e: Tup<>, Void>"),
		nil,
		t,
	)
	interpreter.TestProgram(`{f: {}, e: [], d: "", c: 0, b: false, a: null} type`,
		types.Str{},
		states.StrValue("Obj<a: Null, b: Bool, c: Num, d: Str, e: Tup<>, f: Obj<Void>, Void>"),
		nil,
		t,
	)
}
