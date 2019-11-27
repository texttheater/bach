package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestConditionals(t *testing.T) {
	TestProgram(`if true then 2 else 3 ok`, types.NumType{}, states.NumValue(2), nil, t)
	TestProgram(`for Num def heart Bool as if <3 then true else false ok ok 2 heart`, types.BoolType{}, states.BoolValue(true), nil, t)
	TestProgram(`for Num def heart Bool as if <3 then true else false ok ok 4 heart`, types.BoolType{}, states.BoolValue(false), nil, t)
	TestProgram(`for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 0 -1 expand`, types.NumType{}, states.NumValue(-2), nil, t)
	TestProgram(`for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 1 expand`, types.NumType{}, states.NumValue(2), nil, t)
	TestProgram(`for Num def expand Num as if <0 then -1 elif >0 then +1 else 0 ok ok 0 expand`, types.NumType{}, states.NumValue(0), nil, t)
}
