package tests

import (
	"testing"

	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestLiterals(t *testing.T) {
	TestProgram("1", types.NumType{}, values.NumValue(1), nil, t)
	TestProgram("1 2", types.NumType{}, values.NumValue(2), nil, t)
	TestProgram("1 2 3.5", types.NumType{}, values.NumValue(3.5), nil, t)
}
