package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestSimpleTypes(t *testing.T) {
	TestProgram(`null type`,
		types.Str{},
		states.StrValue("Null"),
		nil,
		t,
	)
	TestProgram(`true type`,
		types.Str{},
		states.StrValue("Bool"),
		nil,
		t,
	)
	TestProgram(`1 type`,
		types.Str{},
		states.StrValue("Num"),
		nil,
		t,
	)
	TestProgram(`"abc" type`,
		types.Str{},
		states.StrValue("Str"),
		nil,
		t,
	)
}
