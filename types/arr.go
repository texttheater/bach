package types

import (
	"fmt"
)

type ArrType struct {
	ElType Type
}

var AnyArrType Type = &ArrType{AnyType{}}

var VoidArrType Type = &ArrType{VoidType{}}

func (t *ArrType) Subsumes(u Type) bool {
	if (VoidType{}).Subsumes(u) {
		return true
	}
	switch u := u.(type) {
	case *ArrType:
		return t.ElType.Subsumes(u.ElType)
	case *NearrType:
		if !t.ElType.Subsumes(u.HeadType) {
			return false
		}
		return t.Subsumes(u.TailType)
	case UnionType:
		for _, disjunct := range u {
			if !t.Subsumes(disjunct) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (t *ArrType) String() string {
	if (VoidType{}).Subsumes(t.ElType) {
		return "Tup<>"
	}
	return fmt.Sprintf("Arr<%s>", t.ElType)
}

func (t *ArrType) ElementType() Type {
	return t.ElType
}
