package types_test

import (
	"testing"

	"github.com/texttheater/bach/types"
)

func TestSubsumption1(t *testing.T) {
	t1 := &types.Arr{types.Num{}}
	t2 := types.NewTup([]types.Type{
		types.Num{},
		types.Num{},
		types.Num{},
	})
	if !t1.Subsumes(t2) {
		t.Logf("%s should subsume %s", t1, t2)
		t.Fail()
	}
}

func TestSubsumption2(t *testing.T) {
	t1 := &types.Arr{types.Void{}}
	t2 := &types.Arr{types.Void{}}
	if !t1.Subsumes(t2) {
		t.Logf("%s should subsume %s", t1, t2)
		t.Fail()
	}
}

func TestSubsumption3(t *testing.T) {
	t1 := &types.Nearr{types.Any{}, &types.Nearr{types.Any{}, &types.Nearr{types.Any{}, types.AnyArr}}}
	t2 := types.NewTup([]types.Type{types.Num{}, types.Num{}})
	if t1.Subsumes(t2) {
		t.Logf("%s should not subsume %s", t1, t2)
		t.Fail()
	}
}

func TestSubsumption4(t *testing.T) {
	a1 := types.Obj{
		map[string]types.Type{"a": types.Num{}},
		types.Void{},
	}
	b2 := types.Obj{
		map[string]types.Type{"b": types.Num{}},
		types.Void{},
	}
	t1 := &types.Nearr{
		a1,
		&types.Nearr{
			b2,
			types.VoidArr,
		},
	}
	gotType := t1.ElementType()
	wantType := types.NewUnion(a1, b2)
	if !gotType.Subsumes(wantType) {
		t.Logf("%s should subsume %s", gotType, wantType)
		t.Fail()
	}
	if !wantType.Subsumes(gotType) {
		t.Logf("%s should subsume %s", wantType, gotType)
		t.Fail()
	}
	gotInt, gotCmp := gotType.Partition(a1)
	wantInt, wantCmp := a1, b2
	if !gotCmp.Subsumes(wantCmp) {
		t.Logf("%s should subsume %s", gotCmp, wantCmp)
		t.Fail()
	}
	if !wantCmp.Subsumes(gotCmp) {
		t.Logf("%s should subsume %s", wantCmp, gotCmp)
		t.Fail()
	}
	if !gotInt.Subsumes(wantInt) {
		t.Logf("%s should subsume %s", gotInt, wantInt)
		t.Fail()
	}
	if !wantInt.Subsumes(gotInt) {
		t.Logf("%s should subsume %s", wantInt, gotInt)
		t.Fail()
	}
}

func TestSubsumption5(t *testing.T) {
	general := types.Obj{
		Props: map[string]types.Type{"yes": types.Num{}},
		Rest:  types.Any{},
	}
	specific := types.Obj{
		Props: map[string]types.Type{"yes": types.Num{}},
		Rest:  types.Void{},
	}
	if !general.Subsumes(specific) {
		t.Logf("%s should subsume %s", general, specific)
		t.Fail()
	}
}
