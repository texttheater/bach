package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestSimpleTypes(t *testing.T) {
	interpreter.TestProgram(`null type`,
		types.Str{},
		states.StrValue("Null"),
		nil,
		t,
	)
	interpreter.TestProgram(`true type`,
		types.Str{},
		states.StrValue("Bool"),
		nil,
		t,
	)
	interpreter.TestProgram(`1 type`,
		types.Str{},
		states.StrValue("Num"),
		nil,
		t,
	)
	interpreter.TestProgram(`"abc" type`,
		types.Str{},
		states.StrValue("Str"),
		nil,
		t,
	)
}
