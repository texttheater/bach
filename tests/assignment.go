package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestAssignment(t *testing.T) {
	TestProgram("1 +1 =a 3 *2 +a", types.Num{}, states.NumValue(8), nil, t)
	TestProgram("1 +1 ==2 =p 1 +1 ==1 =q p ==q not", types.Bool{}, states.BoolValue(true), nil, t)
}
