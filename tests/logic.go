package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestLogic(t *testing.T) {
	TestProgram("true", types.BoolType{}, states.BoolValue(true), nil, t)
	TestProgram("false", types.BoolType{}, states.BoolValue(false), nil, t)
	TestProgram("true and(true)", types.BoolType{}, states.BoolValue(true), nil, t)
	TestProgram("true and(false)", types.BoolType{}, states.BoolValue(false), nil, t)
	TestProgram("false and(false)", types.BoolType{}, states.BoolValue(false), nil, t)
	TestProgram("false and(true)", types.BoolType{}, states.BoolValue(false), nil, t)
	TestProgram("true or(true)", types.BoolType{}, states.BoolValue(true), nil, t)
	TestProgram("true or(false)", types.BoolType{}, states.BoolValue(true), nil, t)
	TestProgram("false or(false)", types.BoolType{}, states.BoolValue(false), nil, t)
	TestProgram("false or(true)", types.BoolType{}, states.BoolValue(true), nil, t)
	TestProgram("false not", types.BoolType{}, states.BoolValue(true), nil, t)
	TestProgram("true not", types.BoolType{}, states.BoolValue(false), nil, t)
	TestProgram("true ==true", types.BoolType{}, states.BoolValue(true), nil, t)
	TestProgram("true ==false", types.BoolType{}, states.BoolValue(false), nil, t)
	TestProgram("false ==false", types.BoolType{}, states.BoolValue(true), nil, t)
	TestProgram("false ==true", types.BoolType{}, states.BoolValue(false), nil, t)
	TestProgram("1 +1 ==2 and(2 +2 ==5 not)", types.BoolType{}, states.BoolValue(true), nil, t)
}
