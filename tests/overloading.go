package tests

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func OverloadingTestCases() []TestCase {
	return []TestCase{
		{`for Bool def f Num as 1 ok for Num def f Num as 2 ok true f`, types.NumType{}, values.NumValue(1), nil},
		{`for Bool def f Num as 1 ok for Num def f Num as 2 ok 1 f`, types.NumType{}, values.NumValue(2), nil},
	}
}
