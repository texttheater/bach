package tests_test

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/tests"
	"github.com/texttheater/bach/types"
)

func TestAnyType(t *testing.T) {
	tests.TestProgram(
		`for Any def f Any as null ok f type`,
		types.Str{},
		states.StrValue("Any"),
		nil,
		t,
	)
}
