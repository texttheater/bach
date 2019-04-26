package tests

import (
	"math"

	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func MathTestCases() []TestCase {
	return []TestCase{
		{"1 +1", types.NumType{}, values.NumValue(2), nil},
		{"1 +2 *3", types.NumType{}, values.NumValue(9), nil},
		{"1 +(2 *3)", types.NumType{}, values.NumValue(7), nil},
		{"1 /0", types.NumType{}, values.NumValue(math.Inf(1)), nil},
		{"0 -1 *2", types.NumType{}, values.NumValue(-2), nil},
		{"15 %7", types.NumType{}, values.NumValue(1), nil},
		{"2 >3", types.BoolType{}, values.BoolValue(false), nil},
		{"2 <3", types.BoolType{}, values.BoolValue(true), nil},
		{"3 >2", types.BoolType{}, values.BoolValue(true), nil},
		{"3 <2", types.BoolType{}, values.BoolValue(false), nil},
		{"1 +1 ==2", types.BoolType{}, values.BoolValue(true), nil},
		{"1 +1 >=2", types.BoolType{}, values.BoolValue(true), nil},
		{"1 +1 <=2", types.BoolType{}, values.BoolValue(true), nil},
	}
}
