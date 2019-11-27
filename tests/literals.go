package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestLiterals(t *testing.T) {
	TestProgram("1", types.NumType{}, states.NumValue(1), nil, t)
	TestProgram("1 2", types.NumType{}, states.NumValue(2), nil, t)
	TestProgram("1 2 3.5", types.NumType{}, states.NumValue(3.5), nil, t)
}
