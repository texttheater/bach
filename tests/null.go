package tests

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func NullTestCases() []TestCase {
	return []TestCase{
		{"1 null", types.NullType{}, &values.NullValue{}, nil},
	}
}
