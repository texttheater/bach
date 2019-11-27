package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestArrayTypes(t *testing.T) {
	TestProgram(`[] type`, types.StrType{}, states.StrValue("Tup<>"), nil, t)
	TestProgram(`["dog", "cat"] type`, types.StrType{}, states.StrValue("Tup<Str, Str>"), nil, t)
	TestProgram(`["dog", 1] type`, types.StrType{}, states.StrValue("Tup<Str, Num>"), nil, t)
	TestProgram(`["dog", 1, {}] type`, types.StrType{}, states.StrValue("Tup<Str, Num, Obj<>>"), nil, t)
	TestProgram(`["dog", 1, {}, 2] type`, types.StrType{}, states.StrValue("Tup<Str, Num, Obj<>, Num>"), nil, t)
}
