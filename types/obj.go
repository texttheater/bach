package types

import (
	"bytes"
	"sort"
)

type ObjType struct {
	Props       []string
	PropTypeMap map[string]Type
}

var AnyObjType = ObjType{
	Props:       []string{},
	PropTypeMap: make(map[string]Type),
}

func NewObjType(propTypeMap map[string]Type) Type {
	props := make([]string, len(propTypeMap))
	i := 0
	for k := range propTypeMap {
		props[i] = k
		i++
	}
	sort.Strings(props)
	return ObjType{
		Props:       props,
		PropTypeMap: propTypeMap,
	}
}

func (t ObjType) Subsumes(u Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case ObjType:
		for k, v1 := range t.PropTypeMap {
			v2, ok := u.PropTypeMap[k]
			if !ok {
				return false
			}
			if !v1.Subsumes(v2) {
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

func (t ObjType) Bind(u Type, bindings map[string]Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case ObjType:
		for k, v1 := range t.PropTypeMap {
			v2, ok := u.PropTypeMap[k]
			if !ok {
				return false
			}
			if !v1.Bind(v2, bindings) {
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

func (t ObjType) Instantiate(bindings map[string]Type) Type {
	propTypeMap := make(map[string]Type)
	for p, t := range t.PropTypeMap {
		propTypeMap[p] = t.Instantiate(bindings)
	}
	return ObjType{
		Props:       t.Props,
		PropTypeMap: propTypeMap,
	}
}

func (t ObjType) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case ObjType:
		propTypeMap := make(map[string]Type)
		allSubsumed := true
		for k, v1 := range t.PropTypeMap {
			if v2, ok := u.PropTypeMap[k]; ok {
				intersection, _ := v1.Partition(v2)
				if (VoidType{}).Subsumes(v1) {
					return VoidType{}, t
				}
				if allSubsumed && !intersection.Subsumes(v1) {
					allSubsumed = false
				}
				v1 = intersection
			}
			propTypeMap[k] = v1
		}
		for k, v2 := range u.PropTypeMap {
			if _, ok := propTypeMap[k]; ok {
				continue
			}
			allSubsumed = false
			propTypeMap[k] = v2
		}
		if allSubsumed {
			return NewObjType(propTypeMap), VoidType{}
		}
		return NewObjType(propTypeMap), t
	case UnionType:
		return u.inversePartition(t)
	case AnyType:
		return t, VoidType{}
	default:
		return VoidType{}, t
	}
}

func (t ObjType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("Obj<")
	if len(t.Props) > 0 {
		buffer.WriteString(t.Props[0])
		buffer.WriteString(": ")
		buffer.WriteString(t.PropTypeMap[t.Props[0]].String())
		for _, prop := range t.Props[1:] {
			typ := t.PropTypeMap[prop]
			buffer.WriteString(", ")
			buffer.WriteString(prop)
			buffer.WriteString(": ")
			buffer.WriteString(typ.String())
		}
	}
	buffer.WriteString(">")
	return buffer.String()
}

func (t ObjType) ElementType() Type {
	panic(t.String() + " is not a sequence type")
}
