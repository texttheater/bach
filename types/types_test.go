package types_test

import (
	"testing"

	"github.com/texttheater/bach/types"
)

func TestSubsumption(t *testing.T) {
	t1 := &types.ArrType{types.NumType{}}
	t2 := types.TupType([]types.Type{
		types.NumType{},
		types.NumType{},
		types.NumType{},
	})
	if !t1.Subsumes(t2) {
		t.Logf("%s should subsume %s", t1, t2)
		t.Fail()
	}
	t1 = &types.ArrType{types.VoidType{}}
	t2 = &types.ArrType{types.VoidType{}}
	if !t1.Subsumes(t2) {
		t.Logf("%s should subsume %s", t1, t2)
		t.Fail()
	}
}
