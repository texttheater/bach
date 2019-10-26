package types

import (
	"bytes"
)

type TupType []Type

func (t TupType) Subsumes(u Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case TupType:
		if len(t) != len(u) {
			return false
		}
		for i := range t {
			if !t[i].Subsumes(u[i]) {
				return false
			}
		}
		return true
	case UnionType:
		return u.inverseSubsumes(t)
	default:
		return false
	}
}

func (t TupType) Bind(u Type, bindings map[string]Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case TupType:
		if len(t) != len(u) {
			return false
		}
		for i := range t {
			if !t[i].Bind(u[i], bindings) {
				return false
			}
		}
		return true
	case UnionType:
		return u.inverseBind(t, bindings)
	default:
		return false
	}
}

func (t TupType) Instantiate(bindings map[string]Type) Type {
	elementTypes := make([]Type, len(t))
	for i, elementType := range t {
		elementTypes[i] = elementType.Instantiate(bindings)
	}
	return TupType(elementTypes)
}

func (t TupType) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case TupType:
		if len(t) != len(u) {
			return VoidType{}, t
		}
		elTypes := make([]Type, len(t))
		allSubsumed := true
		for i := range t {
			intersection, _ := t[i].Partition(u[i])
			if (VoidType{}).Subsumes(intersection) {
				return VoidType{}, t
			}
			allSubsumed = allSubsumed && intersection.Subsumes(t[i])
			elTypes[i] = intersection
		}
		if allSubsumed {
			return TupType(elTypes), VoidType{}
		}
		return TupType(elTypes), t
	case *ArrType:
		elTypes := make([]Type, len(t))
		allSubsumed := true
		for i, elType := range t {
			intersection, _ := elType.Partition(u.ElType)
			if (VoidType{}).Subsumes(intersection) {
				return VoidType{}, t
			}
			allSubsumed = allSubsumed && intersection.Subsumes(elType)
			elTypes[i] = intersection
		}
		if allSubsumed {
			return TupType(elTypes), VoidType{}
		}
		return TupType(elTypes), t
	case *SeqType:
		elTypes := make([]Type, len(t))
		allSubsumed := true
		for i, elType := range t {
			intersection, _ := elType.Partition(u.ElType)
			if (VoidType{}).Subsumes(intersection) {
				return VoidType{}, t
			}
			allSubsumed = allSubsumed && intersection.Subsumes(elType)
			elTypes[i] = intersection
		}
		if allSubsumed {
			return TupType(elTypes), VoidType{}
		}
		return TupType(elTypes), t
	case UnionType:
		return u.inversePartition(t)
	case AnyType:
		return t, VoidType{}
	default:
		return VoidType{}, t
	}
}

func (t TupType) ElementType() Type {
	var elType Type = VoidType{}
	for _, el := range t {
		elType = Union(elType, el)
	}
	return elType
}

func (t TupType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("Tup<")
	if len(t) > 0 {
		buffer.WriteString(t[0].String())
		for _, el := range t[1:] {
			buffer.WriteString(", ")
			buffer.WriteString(el.String())
		}
	}
	buffer.WriteString(">")
	return buffer.String()
}
