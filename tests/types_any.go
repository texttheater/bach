package tests

import (
	"testing"

	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestAnyType(t *testing.T) {
	TestProgram(`for Any def f Any as null ok f type`, types.StrType{}, values.StrValue("Any"), nil, t)
}
