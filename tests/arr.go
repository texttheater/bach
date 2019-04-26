package tests

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func ArrayTestCases() []TestCase {
	return []TestCase{
		{`[]`, types.VoidArrType, values.ArrValue([]values.Value{}), nil},
		{`[1]`, types.TupType(types.NumType{}), values.ArrValue([]values.Value{values.NumValue(1)}), nil},
		{`[1, 2, 3]`, types.TupType(types.NumType{}, types.NumType{}, types.NumType{}), values.ArrValue([]values.Value{values.NumValue(1), values.NumValue(2), values.NumValue(3)}), nil},
		{`[1, "a"]`, types.TupType(types.NumType{}, types.StrType{}), values.ArrValue([]values.Value{values.NumValue(1), values.StrValue("a")}), nil},
		{`[[1, 2], ["a", "b"]]`, types.TupType(types.TupType(types.NumType{}, types.NumType{}), types.TupType(types.StrType{}, types.StrType{})), values.ArrValue([]values.Value{values.ArrValue([]values.Value{values.NumValue(1), values.NumValue(2)}), values.ArrValue([]values.Value{values.StrValue("a"), values.StrValue("b")})}), nil},
	}
}
