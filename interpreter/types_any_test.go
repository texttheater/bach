package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestAnyType(t *testing.T) {
	interpreter.TestProgram(
		`for Any def f Any as null ok f type`,
		types.Str{},
		states.StrValue("Any"),
		nil,
		t,
	)
}
