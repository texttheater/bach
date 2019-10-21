package tests

import (
	"testing"

	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestArrays(t *testing.T) {
	TestProgram(`[]`, types.TupType([]types.Type{}), values.NewArrValue([]values.Value{}), nil, t)
	TestProgram(`[1]`, types.TupType([]types.Type{types.NumType{}}), values.NewArrValue([]values.Value{values.NumValue(1)}), nil, t)
	TestProgram(`[1, 2, 3]`, types.TupType([]types.Type{types.NumType{}, types.NumType{}, types.NumType{}}), values.NewArrValue([]values.Value{values.NumValue(1), values.NumValue(2), values.NumValue(3)}), nil, t)
	TestProgram(`[1, "a"]`, types.TupType([]types.Type{types.NumType{}, types.StrType{}}), values.NewArrValue([]values.Value{values.NumValue(1), values.StrValue("a")}), nil, t)
	TestProgram(`[[1, 2], ["a", "b"]]`, types.TupType([]types.Type{types.TupType([]types.Type{types.NumType{}, types.NumType{}}), types.TupType([]types.Type{types.StrType{}, types.StrType{}})}), values.NewArrValue([]values.Value{values.NewArrValue([]values.Value{values.NumValue(1), values.NumValue(2)}), values.NewArrValue([]values.Value{values.StrValue("a"), values.StrValue("b")})}), nil, t)
}
