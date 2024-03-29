package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestLiterals(t *testing.T) {
	interpreter.TestProgram("1", types.Num{}, states.NumValue(1), nil, t)
	interpreter.TestProgram("1 2", types.Num{}, states.NumValue(2), nil, t)
	interpreter.TestProgram("1 2 3.5", types.Num{}, states.NumValue(3.5), nil, t)
}
