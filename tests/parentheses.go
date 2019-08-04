package tests

import (
	"testing"

	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestParentheses(t *testing.T) {
	TestProgram("1 +2 *3", types.NumType{}, values.NumValue(9), nil, t)
	TestProgram("1 +(2 *3)", types.NumType{}, values.NumValue(7), nil, t)
}
