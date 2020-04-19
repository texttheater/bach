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
	case ObjType:
		for _, valueType := range u.PropTypeMap {
			if !t.ValueType.Subsumes(valueType) {
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

func (t MapType) Bind(u Type, bindings map[string]Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case MapType:
		return t.ValueType.Bind(u.ValueType, bindings)
	case ObjType:
		for _, valueType := range u.PropTypeMap {
			if !t.ValueType.Bind(valueType, bindings) {
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

func (t MapType) Instantiate(bindings map[string]Type) Type {
	return MapType{
		ValueType: t.ValueType.Instantiate(bindings),
	}
}

func (t MapType) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case ObjType:
		propTypeMap := make(map[string]Type)
		for key, valueType := range u.PropTypeMap {
			i, _ := t.ValueType.Partition(valueType)
			if (VoidType{}).Subsumes(i) {
				return VoidType{}, t
			}
			propTypeMap[key] = i
		}
		return NewObjType(propTypeMap), t
	case MapType:
		i, r := t.ValueType.Partition(u.ValueType)
		return MapType{i}, MapType{r}
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
