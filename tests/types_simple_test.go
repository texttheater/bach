package tests_test

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/tests"
	"github.com/texttheater/bach/types"
)

func TestSimpleTypes(t *testing.T) {
	tests.TestProgram(`null type`,
		types.Str{},
		states.StrValue("Null"),
		nil,
		t,
	)
	tests.TestProgram(`true type`,
		types.Str{},
		states.StrValue("Bool"),
		nil,
		t,
	)
	tests.TestProgram(`1 type`,
		types.Str{},
		states.StrValue("Num"),
		nil,
		t,
	)
	tests.TestProgram(`"abc" type`,
		types.Str{},
		states.StrValue("Str"),
		nil,
		t,
	)
}
