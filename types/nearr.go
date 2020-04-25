package types

import (
	"bytes"
)

func TupType(elementTypes []Type) Type {
	return NewNearrType(elementTypes, &ArrType{VoidType{}})
}

func NewNearrType(elementTypes []Type, restType Type) Type {
	var t Type = restType
	for i := len(elementTypes) - 1; i >= 0; i-- {
		t = &NearrType{
			HeadType: elementTypes[i],
			TailType: t,
		}
	}
	return t
}

type NearrType struct {
	HeadType Type
	TailType Type
}

func (t *NearrType) Subsumes(u Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case *NearrType:
		return t.HeadType.Subsumes(u.HeadType) && t.TailType.Subsumes(u.TailType)
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
		return t.HeadType.Bind(u.HeadType, bindings) && t.TailType.Bind(u.TailType, bindings)
	case UnionType:
		return u.inverseBind(t, bindings)
	default:
		return false
	}
}

func (t *NearrType) Instantiate(bindings map[string]Type) Type {
	return &NearrType{
		HeadType: t.HeadType.Instantiate(bindings),
		TailType: t.TailType.Instantiate(bindings),
	}
}

func (t *NearrType) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case *NearrType:
		headIntersection, _ := t.HeadType.Partition(u.HeadType)
		if (VoidType{}).Subsumes(headIntersection) {
			return VoidType{}, t
		}
		tailIntersection, _ := t.TailType.Partition(u.TailType)
		if (VoidType{}).Subsumes(tailIntersection) {
			return VoidType{}, t
		}
		intersection := &NearrType{
			HeadType: headIntersection,
			TailType: tailIntersection,
		}
		if intersection.Subsumes(t) {
			return intersection, VoidType{}
		}
		return intersection, t
	case *ArrType:
		headIntersection, _ := t.HeadType.Partition(u.ElType)
		if (VoidType{}).Subsumes(headIntersection) {
			return VoidType{}, t
		}
		tailIntersection, _ := t.TailType.Partition(u)
		if (VoidType{}).Subsumes(tailIntersection) {
			return VoidType{}, t
		}
		intersection := &NearrType{
			HeadType: headIntersection,
			TailType: tailIntersection,
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
	return Union(t.HeadType, t.TailType.ElementType())
}

func (t *NearrType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("Tup<")
	buffer.WriteString(t.HeadType.String())
	tail := t.TailType
Loop:
	for {
		switch t := tail.(type) {
		case *NearrType:
			buffer.WriteString(", ")
			buffer.WriteString(t.HeadType.String())
			tail = t.TailType
		case *ArrType:
			if !(VoidType{}).Subsumes(t.ElType) {
				buffer.WriteString(", ")
				buffer.WriteString(t.ElType.String())
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
