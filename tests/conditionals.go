package tests

import (
	"testing"

	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestConditionals(t *testing.T) {
	TestProgram(`if true then 2 else 3 ok`, types.NumType{}, values.NumValue(2), nil, t)
	TestProgram(`for Num def heart Bool as if <3 then true else false ok ok 2 heart`, types.BoolType{}, values.BoolValue(true), nil, t)
	TestProgram(`for Num def heart Bool as if <3 then true else false ok ok 4 heart`, types.BoolType{}, values.BoolValue(false), nil, t)
	TestProgram(`for Num def expand Num as if <0 then -1 elseIf >0 then +1 else 0 ok ok 0 -1 expand`, types.NumType{}, values.NumValue(-2), nil, t)
	TestProgram(`for Num def expand Num as if <0 then -1 elseIf >0 then +1 else 0 ok ok 1 expand`, types.NumType{}, values.NumValue(2), nil, t)
	TestProgram(`for Num def expand Num as if <0 then -1 elseIf >0 then +1 else 0 ok ok 0 expand`, types.NumType{}, values.NumValue(0), nil, t)
}
