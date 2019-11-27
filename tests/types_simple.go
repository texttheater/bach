package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestSimpleTypes(t *testing.T) {
	TestProgram(`null type`, types.StrType{}, states.StrValue("Null"), nil, t)
	TestProgram(`true type`, types.StrType{}, states.StrValue("Bool"), nil, t)
	TestProgram(`1 type`, types.StrType{}, states.StrValue("Num"), nil, t)
	TestProgram(`"abc" type`, types.StrType{}, states.StrValue("Str"), nil, t)
}
