package interpreter_test

import (
	"math"
	"testing"

	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestMath(t *testing.T) {
	interpreter.TestProgram("1 +1", types.Num{}, states.NumValue(2), nil, t)
	interpreter.TestProgram("1 +2 *3", types.Num{}, states.NumValue(9), nil, t)
	interpreter.TestProgram("1 +(2 *3)", types.Num{}, states.NumValue(7), nil, t)
	interpreter.TestProgram("1 /0", types.Num{}, states.NumValue(math.Inf(1)), nil, t)
	interpreter.TestProgram("0 -1 *2", types.Num{}, states.NumValue(-2), nil, t)
	interpreter.TestProgram("-1 *2", types.Num{}, states.NumValue(-2), nil, t)
	interpreter.TestProgram("-0 *2", types.Num{}, states.NumValue(math.Copysign(0, -1)), nil, t)
	interpreter.TestProgram("2 **5", types.Num{}, states.NumValue(32), nil, t)
	interpreter.TestProgram("15 %7", types.Num{}, states.NumValue(1), nil, t)
	interpreter.TestProgram("2 >3", types.Bool{}, states.BoolValue(false), nil, t)
	interpreter.TestProgram("2 <3", types.Bool{}, states.BoolValue(true), nil, t)
	interpreter.TestProgram("3 >2", types.Bool{}, states.BoolValue(true), nil, t)
	interpreter.TestProgram("3 <2", types.Bool{}, states.BoolValue(false), nil, t)
	interpreter.TestProgram("1 +1 >=2", types.Bool{}, states.BoolValue(true), nil, t)
	interpreter.TestProgram("1 +1 <=2", types.Bool{}, states.BoolValue(true), nil, t)
}
