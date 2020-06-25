package types_test

import (
	"testing"

	"github.com/texttheater/bach/types"
)

func TestSubsumption(t *testing.T) {
	var t1 types.Type = &types.ArrType{types.NumType{}}
	var t2 types.Type = types.TupType([]types.Type{
		types.NumType{},
		types.NumType{},
		types.NumType{},
	})
	if !t1.Subsumes(t2) {
		t.Logf("%s should subsume %s", t1, t2)
		t.Fail()
	}
	// ---
	t1 = &types.ArrType{types.VoidType{}}
	t2 = &types.ArrType{types.VoidType{}}
	if !t1.Subsumes(t2) {
		t.Logf("%s should subsume %s", t1, t2)
		t.Fail()
	}
	// ---
	t1 = &types.NearrType{types.AnyType{}, &types.NearrType{types.AnyType{}, &types.NearrType{types.AnyType{}, types.AnyArrType}}}
	t2 = types.TupType([]types.Type{types.NumType{}, types.NumType{}})
	if t1.Subsumes(t2) {
		t.Logf("%s should not subsume %s", t1, t2)
		t.Fail()
	}
}
