package tests

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func AssignmentTestCases() []TestCase {
	return []TestCase{
		{"1 +1 =a 3 *2 +a", types.NumType{}, values.NumValue(8), nil},
		{"1 +1 ==2 =p 1 +1 ==1 =q p ==q not", types.BoolType{}, values.BoolValue(true), nil},
	}
}
