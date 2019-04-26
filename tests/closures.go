package tests

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func ClosureTestCases() []TestCase {
	return []TestCase{
		{`1 =a for Any def f Num as a ok f 2 =a f`, types.NumType{}, values.NumValue(1), nil},
	}
}
