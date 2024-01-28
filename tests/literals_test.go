package tests_test

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/tests"
	"github.com/texttheater/bach/types"
)

func TestLiterals(t *testing.T) {
	tests.TestProgram("1", types.Num{}, states.NumValue(1), nil, t)
	tests.TestProgram("1 2", types.Num{}, states.NumValue(2), nil, t)
	tests.TestProgram("1 2 3.5", types.Num{}, states.NumValue(3.5), nil, t)
}
