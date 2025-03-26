package types

import (
	"fmt"
)

type ArrType struct {
	El Type
}

var AnyArrType Type = &ArrType{AnyType{}}

var VoidArrType Type = &ArrType{VoidType{}}

func NewArrType(el Type) *ArrType {
	return &ArrType{
		El: el,
	}
}

func (t *ArrType) Subsumes(u Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case *NearrType:
		return t.El.Subsumes(u.Head) && t.Subsumes(u.Tail)
	case *ArrType:
		return t.El.Subsumes(u.El)
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
		return t.El.Bind(u.ElementType(), bindings)
	case *ArrType:
		return t.El.Bind(u.El, bindings)
	case UnionType:
		return u.inverseBind(t, bindings)
	default:
		return false
	}
}

func (t *ArrType) Instantiate(bindings map[string]Type) Type {
	return &ArrType{t.El.Instantiate(bindings)}
}

func (t *ArrType) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case *NearrType:
		headIntersection, _ := t.El.Partition(u.Head)
		if (VoidType{}).Subsumes(headIntersection) {
			return VoidType{}, t
		}
		tailIntersection, _ := t.Partition(u.Tail)
		if (VoidType{}).Subsumes(tailIntersection) {
			return VoidType{}, t
		}
		return &NearrType{
			Head: headIntersection,
			Tail: tailIntersection,
		}, t
	case *ArrType:
		intersection, _ := t.El.Partition(u.El)
		if intersection.Subsumes(t.El) {
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
	if (VoidType{}).Subsumes(t.El) {
		return "Arr<>"
	}
	return fmt.Sprintf("Arr<%s...>", t.El)
}

func (t *ArrType) ElementType() Type {
	return t.El
}
