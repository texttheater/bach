package tests

import (
	"testing"

	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestSimpleTypes(t *testing.T) {
	TestProgram(`null type`, types.StrType{}, values.StrValue("Null"), nil, t)
	TestProgram(`true type`, types.StrType{}, values.StrValue("Bool"), nil, t)
	TestProgram(`1 type`, types.StrType{}, values.StrValue("Num"), nil, t)
	TestProgram(`"abc" type`, types.StrType{}, values.StrValue("Str"), nil, t)
}
