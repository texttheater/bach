package tests

import (
	"testing"

	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestMatching(t *testing.T) {
	TestProgram(`if true then 2 else "two" ok is Num then true else false ok`, types.BoolType{}, values.BoolValue(true), nil, t)
}
