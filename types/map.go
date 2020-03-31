package types

import (
	"bytes"
)

type MapType struct {
	ValueType Type
}

func (t MapType) Subsumes(u Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case MapType:
		return t.ValueType.Subsumes(u.ValueType)
	case UnionType:
		return u.inverseSubsumes(t)
	default:
		return false
	}
}

func (t MapType) Bind(u Type, bindings map[string]Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case MapType:
		return t.ValueType.Bind(u.ValueType, bindings)
	case UnionType:
		return u.inverseBind(t, bindings)
	default:
		return false
	}
}

func (t MapType) Instantiate(bindings map[string]Type) Type {
	return MapType{
		ValueType: t.ValueType.Instantiate(bindings),
	}
}

func (t MapType) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case MapType:
		i, r := t.ValueType.Partition(u.ValueType)
		return MapType{i}, MapType{r}
	case ObjType:
		for _, v := range u.PropTypeMap {
			if !t.ValueType.Subsumes(v) && !v.Subsumes(t.ValueType) {
				return VoidType{}, t
			}
		}
		return t, t
	case UnionType:
		return u.inversePartition(t)
	case AnyType:
		return t, VoidType{}
	default:
		return VoidType{}, t
	}
}

func (t MapType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("Map<")
	buffer.WriteString(t.ValueType.String())
	buffer.WriteString(">")
	return buffer.String()
}

func (t MapType) ElementType() Type {
	panic(t.String() + " is not a sequence type")
}
