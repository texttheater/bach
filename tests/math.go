package tests

import (
	"math"
	"testing"

	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestMath(t *testing.T) {
	TestProgram("1 +1", types.NumType{}, values.NumValue(2), nil, t)
	TestProgram("1 +2 *3", types.NumType{}, values.NumValue(9), nil, t)
	TestProgram("1 +(2 *3)", types.NumType{}, values.NumValue(7), nil, t)
	TestProgram("1 /0", types.NumType{}, values.NumValue(math.Inf(1)), nil, t)
	TestProgram("0 -1 *2", types.NumType{}, values.NumValue(-2), nil, t)
	TestProgram("15 %7", types.NumType{}, values.NumValue(1), nil, t)
	TestProgram("2 >3", types.BoolType{}, values.BoolValue(false), nil, t)
	TestProgram("2 <3", types.BoolType{}, values.BoolValue(true), nil, t)
	TestProgram("3 >2", types.BoolType{}, values.BoolValue(true), nil, t)
	TestProgram("3 <2", types.BoolType{}, values.BoolValue(false), nil, t)
	TestProgram("1 +1 ==2", types.BoolType{}, values.BoolValue(true), nil, t)
	TestProgram("1 +1 >=2", types.BoolType{}, values.BoolValue(true), nil, t)
	TestProgram("1 +1 <=2", types.BoolType{}, values.BoolValue(true), nil, t)
}
