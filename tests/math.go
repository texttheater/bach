package tests

import (
	"math"
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestMath(t *testing.T) {
	TestProgram("1 +1", types.NumType{}, states.NumValue(2), nil, t)
	TestProgram("1 +2 *3", types.NumType{}, states.NumValue(9), nil, t)
	TestProgram("1 +(2 *3)", types.NumType{}, states.NumValue(7), nil, t)
	TestProgram("1 /0", types.NumType{}, states.NumValue(math.Inf(1)), nil, t)
	TestProgram("0 -1 *2", types.NumType{}, states.NumValue(-2), nil, t)
	TestProgram("15 %7", types.NumType{}, states.NumValue(1), nil, t)
	TestProgram("2 >3", types.BoolType{}, states.BoolValue(false), nil, t)
	TestProgram("2 <3", types.BoolType{}, states.BoolValue(true), nil, t)
	TestProgram("3 >2", types.BoolType{}, states.BoolValue(true), nil, t)
	TestProgram("3 <2", types.BoolType{}, states.BoolValue(false), nil, t)
	TestProgram("1 +1 ==2", types.BoolType{}, states.BoolValue(true), nil, t)
	TestProgram("1 +1 >=2", types.BoolType{}, states.BoolValue(true), nil, t)
	TestProgram("1 +1 <=2", types.BoolType{}, states.BoolValue(true), nil, t)
}
