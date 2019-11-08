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
	switch u := u.(type) {
	case VoidType:
		return true
	case *NearrType:
		return t.ElType.Subsumes(u.HeadType) && t.Subsumes(u.TailType)
	case *ArrType:
		return t.ElType.Subsumes(u.ElType)
	case UnionType:
		return u.inverseSubsumes(t)
	default:
		return false
	}
}

func (t *ArrType) Bind(u Type, bindings map[string]Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case *NearrType:
		return t.ElType.Bind(u.HeadType, bindings) && t.Bind(u.TailType, bindings)
	case *ArrType:
		return t.ElType.Bind(u.ElType, bindings)
	case UnionType:
		return u.inverseBind(t, bindings)
	default:
		return false
	}
}

func (t *ArrType) Instantiate(bindings map[string]Type) Type {
	return t
}

func (t *ArrType) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case *NearrType:
		headIntersection, _ := t.ElType.Partition(u.HeadType)
		if (VoidType{}).Subsumes(headIntersection) {
			return VoidType{}, t
		}
		tailIntersection, _ := t.Partition(u.TailType)
		if (VoidType{}).Subsumes(tailIntersection) {
			return VoidType{}, t
		}
		return &NearrType{
			HeadType: headIntersection,
			TailType: tailIntersection,
		}, t
	case *ArrType:
		intersection, _ := t.ElType.Partition(u.ElType)
		if intersection.Subsumes(t.ElType) {
			return &ArrType{intersection}, VoidType{}
		}
		return &ArrType{intersection}, t
	case UnionType:
		return u.inversePartition(t)
	case AnyType:
		return t, VoidType{}
	default:
		return VoidType{}, t
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
