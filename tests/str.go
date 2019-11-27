package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestStrings(t *testing.T) {
	TestProgram(`"abc"`, types.StrType{}, states.StrValue("abc"), nil, t)
	TestProgram(`"\"\\abc\""`, types.StrType{}, states.StrValue(`"\abc"`), nil, t)
	TestProgram(`1 "abc"`, types.StrType{}, states.StrValue("abc"), nil, t)
}
