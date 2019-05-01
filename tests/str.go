package tests

import (
	"testing"

	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestStrings(t *testing.T) {
	TestProgram(`"abc"`, types.StrType{}, values.StrValue("abc"), nil, t)
	TestProgram(`"\"\\abc\""`, types.StrType{}, values.StrValue(`"\abc"`), nil, t)
	TestProgram(`1 "abc"`, types.StrType{}, values.StrValue("abc"), nil, t)
}
