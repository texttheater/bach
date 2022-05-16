package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestLogic(t *testing.T) {
	TestProgram("true", types.Bool{}, states.BoolValue(true), nil, t)
	TestProgram("false", types.Bool{}, states.BoolValue(false), nil, t)
	TestProgram("true and(true)", types.Bool{}, states.BoolValue(true), nil, t)
	TestProgram("true and(false)", types.Bool{}, states.BoolValue(false), nil, t)
	TestProgram("false and(false)", types.Bool{}, states.BoolValue(false), nil, t)
	TestProgram("false and(true)", types.Bool{}, states.BoolValue(false), nil, t)
	TestProgram("true or(true)", types.Bool{}, states.BoolValue(true), nil, t)
	TestProgram("true or(false)", types.Bool{}, states.BoolValue(true), nil, t)
	TestProgram("false or(false)", types.Bool{}, states.BoolValue(false), nil, t)
	TestProgram("false or(true)", types.Bool{}, states.BoolValue(true), nil, t)
	TestProgram("false not", types.Bool{}, states.BoolValue(true), nil, t)
	TestProgram("true not", types.Bool{}, states.BoolValue(false), nil, t)
	TestProgram("1 +1 ==2 and(2 +2 ==5 not)", types.Bool{}, states.BoolValue(true), nil, t)
}
