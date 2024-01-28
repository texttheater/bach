package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestArrayTypes(t *testing.T) {
	interpreter.TestProgram(
		`[] type`,
		types.Str{},
		states.StrValue("Tup<>"),
		nil,
		t,
	)
	interpreter.TestProgram(
		`["dog", "cat"] type`,
		types.Str{},
		states.StrValue("Tup<Str, Str>"),
		nil,
		t,
	)
	interpreter.TestProgram(
		`["dog", 1] type`,
		types.Str{},
		states.StrValue("Tup<Str, Num>"),
		nil,
		t,
	)
	interpreter.TestProgram(
		`["dog", 1, {}] type`,
		types.Str{},
		states.StrValue("Tup<Str, Num, Obj<Void>>"),
		nil,
		t,
	)
	interpreter.TestProgram(
		`["dog", 1, {}, 2] type`,
		types.Str{},
		states.StrValue("Tup<Str, Num, Obj<Void>, Num>"),
		nil,
		t,
	)
}
