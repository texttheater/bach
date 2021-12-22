package types_test

import (
	"testing"

	"github.com/texttheater/bach/types"
)

func TestSubsumption(t *testing.T) {
	var t1 types.Type = &types.Arr{types.Num{}}
	var t2 types.Type = types.NewTup([]types.Type{
		types.Num{},
		types.Num{},
		types.Num{},
	})
	if !t1.Subsumes(t2) {
		t.Logf("%s should subsume %s", t1, t2)
		t.Fail()
	}
	// ---
	t1 = &types.Arr{types.Void{}}
	t2 = &types.Arr{types.Void{}}
	if !t1.Subsumes(t2) {
		t.Logf("%s should subsume %s", t1, t2)
		t.Fail()
	}
	// ---
	t1 = &types.Nearr{types.Any{}, &types.Nearr{types.Any{}, &types.Nearr{types.Any{}, types.AnyArr}}}
	t2 = types.NewTup([]types.Type{types.Num{}, types.Num{}})
	if t1.Subsumes(t2) {
		t.Logf("%s should not subsume %s", t1, t2)
		t.Fail()
	}
}
