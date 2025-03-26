package types

import (
	"bytes"
)

func NewTup(elementTypes []Type) Type {
	return NewNearr(elementTypes, &ArrType{VoidType{}})
}

func NewNearr(elementTypes []Type, restType Type) Type {
	var t Type = restType
	for i := len(elementTypes) - 1; i >= 0; i-- {
		t = &NearrType{
			Head: elementTypes[i],
			Tail: t,
		}
	}
	return t
}

type NearrType struct {
	Head Type
	Tail Type
}

func (t *NearrType) Subsumes(u Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case *NearrType:
		return t.Head.Subsumes(u.Head) && t.Tail.Subsumes(u.Tail)
	case UnionType:
		return u.inverseSubsumes(t)
	default:
		return false
	}
}

func (t *NearrType) Bind(u Type, bindings map[string]Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case *NearrType:
		return t.Head.Bind(u.Head, bindings) && t.Tail.Bind(u.Tail, bindings)
	case UnionType:
		return u.inverseBind(t, bindings)
	default:
		return false
	}
}

func (t *NearrType) Instantiate(bindings map[string]Type) Type {
	return &NearrType{
		Head: t.Head.Instantiate(bindings),
		Tail: t.Tail.Instantiate(bindings),
	}
}

func (t *NearrType) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case *NearrType:
		headIntersection, _ := t.Head.Partition(u.Head)
		if (VoidType{}).Subsumes(headIntersection) {
			return VoidType{}, t
		}
		tailIntersection, _ := t.Tail.Partition(u.Tail)
		if (VoidType{}).Subsumes(tailIntersection) {
			return VoidType{}, t
		}
		intersection := &NearrType{
			Head: headIntersection,
			Tail: tailIntersection,
		}
		if intersection.Subsumes(t) {
			return intersection, VoidType{}
		}
		return intersection, t
	case *ArrType:
		headIntersection, _ := t.Head.Partition(u.El)
		if (VoidType{}).Subsumes(headIntersection) {
			return VoidType{}, t
		}
		tailIntersection, _ := t.Tail.Partition(u)
		if (VoidType{}).Subsumes(tailIntersection) {
			return VoidType{}, t
		}
		intersection := &NearrType{
			Head: headIntersection,
			Tail: tailIntersection,
		}
		if intersection.Subsumes(t) {
			return intersection, VoidType{}
		}
		return intersection, t
	case UnionType:
		return u.inversePartition(t)
	case AnyType:
		return t, VoidType{}
	default:
		return VoidType{}, t
	}
}

func (t NearrType) ElementType() Type {
	return NewUnionType(t.Head, t.Tail.ElementType())
}

func (t *NearrType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("Arr<")
	buffer.WriteString(t.Head.String())
	tail := t.Tail
Loop:
	for {
		switch t := tail.(type) {
		case *NearrType:
			buffer.WriteString(", ")
			buffer.WriteString(t.Head.String())
			tail = t.Tail
		case *ArrType:
			if !(VoidType{}).Subsumes(t.El) {
				buffer.WriteString(", ")
				buffer.WriteString(t.El.String())
				buffer.WriteString("...")
			}
			break Loop
		default:
			panic("non-array tail")
		}
	}
	buffer.WriteString(">")
	return buffer.String()
}
