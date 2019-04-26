package tests

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func RecursionTestCases() []TestCase {
	return []TestCase{
		{`for Num def fac Num as if ==0 then 1 else =n -1 fac *n ok ok 3 fac`, types.NumType{}, values.NumValue(6), nil},
	}
}
