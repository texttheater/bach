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
	case TupType:
		for _, el := range u {
			if !t.ElType.Subsumes(el) {
				return false
			}
		}
		return true
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
	case TupType:
		for _, el := range u {
			if !t.ElType.Bind(el, bindings) {
				return false
			}
		}
		return true
	case *ArrType:
		return t.ElType.Bind(u.ElType, bindings)
	case UnionType:
		return u.inverseBind(t, bindings)
	default:
		return false
	}
}

func (t *ArrType) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case TupType:
		elTypes := make([]Type, len(u))
		for i, elType := range u {
			intersection, _ := t.ElType.Partition(elType)
			if (VoidType{}).Subsumes(intersection) {
				return VoidType{}, t
			}
			elTypes[i] = intersection
		}
		return TupType(elTypes), t
	case *ArrType:
		intersection, _ := t.ElType.Partition(u.ElType)
		if (VoidType{}).Subsumes(intersection) {
			return VoidType{}, t
		}
		if intersection.Subsumes(t.ElType) {
			return &ArrType{intersection}, VoidType{}
		}
		return &ArrType{intersection}, t
	case *SeqType:
		intersection, _ := t.ElType.Partition(u.ElType)
		if (VoidType{}).Subsumes(intersection) {
			return VoidType{}, t
		}
		if intersection.Subsumes(t.ElType) {
			return &SeqType{intersection}, VoidType{}
		}
		return &SeqType{intersection}, t
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
