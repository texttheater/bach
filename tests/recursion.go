package tests

import (
	"testing"

	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestRecursion(t *testing.T) {
	TestProgram(`for Num def fac Num as if ==0 then 1 else =n -1 fac *n ok ok 3 fac`, types.NumType{}, values.NumValue(6), nil, t)
}
