package tests

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func LiteralTestCases() []TestCase {
	return []TestCase{
		{"1", types.NumType{}, values.NumValue(1), nil},
		{"1 2", types.NumType{}, values.NumValue(2), nil},
		{"1 2 3.5", types.NumType{}, values.NumValue(3.5), nil},
	}
}
