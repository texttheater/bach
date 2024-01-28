package tests_test

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/tests"
	"github.com/texttheater/bach/types"
)

func TestLogic(t *testing.T) {
	tests.TestProgram("true", types.Bool{}, states.BoolValue(true), nil, t)
	tests.TestProgram("false", types.Bool{}, states.BoolValue(false), nil, t)
	tests.TestProgram("true and(true)", types.Bool{}, states.BoolValue(true), nil, t)
	tests.TestProgram("true and(false)", types.Bool{}, states.BoolValue(false), nil, t)
	tests.TestProgram("false and(false)", types.Bool{}, states.BoolValue(false), nil, t)
	tests.TestProgram("false and(true)", types.Bool{}, states.BoolValue(false), nil, t)
	tests.TestProgram("true or(true)", types.Bool{}, states.BoolValue(true), nil, t)
	tests.TestProgram("true or(false)", types.Bool{}, states.BoolValue(true), nil, t)
	tests.TestProgram("false or(false)", types.Bool{}, states.BoolValue(false), nil, t)
	tests.TestProgram("false or(true)", types.Bool{}, states.BoolValue(true), nil, t)
	tests.TestProgram("false not", types.Bool{}, states.BoolValue(true), nil, t)
	tests.TestProgram("true not", types.Bool{}, states.BoolValue(false), nil, t)
	tests.TestProgram("1 +1 ==2 and(2 +2 ==5 not)", types.Bool{}, states.BoolValue(true), nil, t)
}
