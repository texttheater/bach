package types

import (
	"fmt"
)

type SeqType struct {
	ElType Type
}

// AnySeqType represents Seq<Any>, the type of sequences with any type of
// element. It is provided as a package variable for convenience.
var AnySeqType Type = &SeqType{AnyType{}}

func (t *SeqType) Subsumes(u Type) bool {
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
	case *SeqType:
		return t.ElType.Subsumes(u.ElType)
	case UnionType:
		return u.inverseSubsumes(t)
	default:
		return false
	}
}

func (t *SeqType) Bind(u Type, bindings map[string]Type) bool {
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
	case *SeqType:
		return t.ElType.Bind(u.ElType, bindings)
	case UnionType:
		return u.inverseBind(t, bindings)
	default:
		return false
	}
}

func (t *SeqType) Instantiate(bindings map[string]Type) Type {
	return &SeqType{
		ElType: t.ElType.Instantiate(bindings),
	}
}

func (t *SeqType) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case TupType:
		elTypes := make([]Type, len(u))
		for i, el := range u {
			intersection, _ := t.ElType.Partition(el)
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

func (t *SeqType) String() string {
	return fmt.Sprintf("Seq<%v>", t.ElType)
}

func (t *SeqType) ElementType() Type {
	return t.ElType
}
