package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestLogic(t *testing.T) {
	interpreter.TestProgram("true", types.Bool{}, states.BoolValue(true), nil, t)
	interpreter.TestProgram("false", types.Bool{}, states.BoolValue(false), nil, t)
	interpreter.TestProgram("true and(true)", types.Bool{}, states.BoolValue(true), nil, t)
	interpreter.TestProgram("true and(false)", types.Bool{}, states.BoolValue(false), nil, t)
	interpreter.TestProgram("false and(false)", types.Bool{}, states.BoolValue(false), nil, t)
	interpreter.TestProgram("false and(true)", types.Bool{}, states.BoolValue(false), nil, t)
	interpreter.TestProgram("true or(true)", types.Bool{}, states.BoolValue(true), nil, t)
	interpreter.TestProgram("true or(false)", types.Bool{}, states.BoolValue(true), nil, t)
	interpreter.TestProgram("false or(false)", types.Bool{}, states.BoolValue(false), nil, t)
	interpreter.TestProgram("false or(true)", types.Bool{}, states.BoolValue(true), nil, t)
	interpreter.TestProgram("false not", types.Bool{}, states.BoolValue(true), nil, t)
	interpreter.TestProgram("true not", types.Bool{}, states.BoolValue(false), nil, t)
	interpreter.TestProgram("1 +1 ==2 and(2 +2 ==5 not)", types.Bool{}, states.BoolValue(true), nil, t)
}
