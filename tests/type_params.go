package tests

import (
	"testing"

	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestTypeParams(t *testing.T) {
	TestProgram(
		`for <A> A def f Str as type ok [1 f, "a" f]`,
		types.TupType([]types.Type{
			types.StrType{},
			types.StrType{},
		}),
		values.ArrValue([]values.Value{
			values.StrValue("Num"),
			values.StrValue("Str"),
		}),
		nil,
		t,
	)
}
