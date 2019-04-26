package tests

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func LogicTestCases() []TestCase {
	return []TestCase{
		{"true", types.BoolType{}, values.BoolValue(true), nil},
		{"false", types.BoolType{}, values.BoolValue(false), nil},
		{"true and(true)", types.BoolType{}, values.BoolValue(true), nil},
		{"true and(false)", types.BoolType{}, values.BoolValue(false), nil},
		{"false and(false)", types.BoolType{}, values.BoolValue(false), nil},
		{"false and(true)", types.BoolType{}, values.BoolValue(false), nil},
		{"true or(true)", types.BoolType{}, values.BoolValue(true), nil},
		{"true or(false)", types.BoolType{}, values.BoolValue(true), nil},
		{"false or(false)", types.BoolType{}, values.BoolValue(false), nil},
		{"false or(true)", types.BoolType{}, values.BoolValue(true), nil},
		{"false not", types.BoolType{}, values.BoolValue(true), nil},
		{"true not", types.BoolType{}, values.BoolValue(false), nil},
		{"true ==true", types.BoolType{}, values.BoolValue(true), nil},
		{"true ==false", types.BoolType{}, values.BoolValue(false), nil},
		{"false ==false", types.BoolType{}, values.BoolValue(true), nil},
		{"false ==true", types.BoolType{}, values.BoolValue(false), nil},
		{"1 +1 ==2 and(2 +2 ==5 not)", types.BoolType{}, values.BoolValue(true), nil},
	}
}
