package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestTemplateLiterals(t *testing.T) {
	TestProgram(
		"`a{2 +2}`",
		types.StrType{},
		states.StrValue("a4"),
		nil,
		t,
	)
}
