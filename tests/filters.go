package tests

import (
	"testing"

	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestFilters(t *testing.T) {
	TestProgram(
		`["a", 1, "b", 2, "c", 3] eachis Num with %2 >0 elis Str all arr`,
		&types.ArrType{types.Union(types.NumType{}, types.StrType{})},
		values.ArrValue([]values.Value{
			values.StrValue("a"),
			values.NumValue(1),
			values.StrValue("b"),
			values.StrValue("c"),
			values.NumValue(3),
		}),
		nil,
		t,
	)
}
