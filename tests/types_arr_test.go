package tests_test

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/tests"
	"github.com/texttheater/bach/types"
)

func TestArrayTypes(t *testing.T) {
	tests.TestProgram(
		`[] type`,
		types.Str{},
		states.StrValue("Tup<>"),
		nil,
		t,
	)
	tests.TestProgram(
		`["dog", "cat"] type`,
		types.Str{},
		states.StrValue("Tup<Str, Str>"),
		nil,
		t,
	)
	tests.TestProgram(
		`["dog", 1] type`,
		types.Str{},
		states.StrValue("Tup<Str, Num>"),
		nil,
		t,
	)
	tests.TestProgram(
		`["dog", 1, {}] type`,
		types.Str{},
		states.StrValue("Tup<Str, Num, Obj<Void>>"),
		nil,
		t,
	)
	tests.TestProgram(
		`["dog", 1, {}, 2] type`,
		types.Str{},
		states.StrValue("Tup<Str, Num, Obj<Void>, Num>"),
		nil,
		t,
	)
}
