package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestAnyType(t *testing.T) {
	TestProgram(`for Any def f Any as null ok f type`, types.StrType{}, states.StrValue("Any"), nil, t)
}
