package types_test

import (
	"testing"

	"github.com/texttheater/bach/types"
)

func TestSubsumption1(t *testing.T) {
	t1 := &types.ArrType{types.NumType{}}
	t2 := types.NewTup([]types.Type{
		types.NumType{},
		types.NumType{},
		types.NumType{},
	})
	if !t1.Subsumes(t2) {
		t.Logf("%s should subsume %s", t1, t2)
		t.Fail()
	}
}

func TestSubsumption2(t *testing.T) {
	t1 := &types.ArrType{types.VoidType{}}
	t2 := &types.ArrType{types.VoidType{}}
	if !t1.Subsumes(t2) {
		t.Logf("%s should subsume %s", t1, t2)
		t.Fail()
	}
}

func TestSubsumption3(t *testing.T) {
	t1 := &types.NearrType{types.AnyType{}, &types.NearrType{types.AnyType{}, &types.NearrType{types.AnyType{}, types.AnyArrType}}}
	t2 := types.NewTup([]types.Type{types.NumType{}, types.NumType{}})
	if t1.Subsumes(t2) {
		t.Logf("%s should not subsume %s", t1, t2)
		t.Fail()
	}
}

func TestSubsumption4(t *testing.T) {
	a1 := types.ObjType{
		map[string]types.Type{"a": types.NumType{}},
		types.VoidType{},
	}
	b2 := types.ObjType{
		map[string]types.Type{"b": types.NumType{}},
		types.VoidType{},
	}
	t1 := &types.NearrType{
		a1,
		&types.NearrType{
			b2,
			types.VoidArrType,
		},
	}
	gotType := t1.ElementType()
	wantType := types.NewUnionType(a1, b2)
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
	general := types.ObjType{
		Props: map[string]types.Type{"yes": types.NumType{}},
		Rest:  types.AnyType{},
	}
	specific := types.ObjType{
		Props: map[string]types.Type{"yes": types.NumType{}},
		Rest:  types.VoidType{},
	}
	if !general.Subsumes(specific) {
		t.Logf("%s should subsume %s", general, specific)
		t.Fail()
	}
}
