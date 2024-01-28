package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestNull(t *testing.T) {
	interpreter.TestProgram("1 null", types.Null{}, states.NullValue{}, nil, t)
}
