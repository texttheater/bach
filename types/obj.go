package types

import (
	"bytes"
	"sort"
)

type ObjType struct {
	PropTypeMap map[string]Type
	RestType    Type
}

var VoidObjType Type = ObjType{
	PropTypeMap: make(map[string]Type),
	RestType:    VoidType{},
}

var AnyObjType Type = ObjType{
	PropTypeMap: make(map[string]Type),
	RestType:    AnyType{},
}

func (t ObjType) Subsumes(u Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case ObjType:
		for prop, wantType := range t.PropTypeMap {
			gotType, ok := u.PropTypeMap[prop]
			if !ok {
				return false
			}
			if !wantType.Subsumes(gotType) {
				return false
			}
		}
		if !t.RestType.Subsumes(u.RestType) {
			return false
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
		for prop, wantType := range t.PropTypeMap {
			gotType, ok := u.PropTypeMap[prop]
			if !ok {
				return false
			}
			if !wantType.Bind(gotType, bindings) {
				return false
			}
		}
		restType := u.RestType
		for prop, gotType := range u.PropTypeMap {
			if _, ok := t.PropTypeMap[prop]; !ok {
				restType = Union(restType, gotType)
			}
		}
		if !t.RestType.Bind(restType, bindings) {
			return false
		}
		return true
	case UnionType:
		return u.inverseSubsumes(t)
	default:
		return false
	}
}

func (t ObjType) Partition(u Type) (Type, Type) {
	switch u := u.(type) {
	case VoidType:
		return u, t
	case ObjType:
		propTypeMap := make(map[string]Type)
		allSubsumed := true
		for prop, tValueType := range t.PropTypeMap {
			uValueType := u.PropTypeMap[prop]
			if uValueType != nil {
				i, c := tValueType.Partition(uValueType)
				if (VoidType{}).Subsumes(i) {
					return VoidType{}, t
				}
				propTypeMap[prop] = i
				allSubsumed = allSubsumed && (VoidType{}).Subsumes(c)
			} else {
				i, c := tValueType.Partition(u.RestType)
				if (VoidType{}).Subsumes(i) {
					return VoidType{}, t
				}
				propTypeMap[prop] = i
				allSubsumed = allSubsumed && (VoidType{}).Subsumes(c)
			}
		}
		for prop, uValueType := range u.PropTypeMap {
			if t.PropTypeMap[prop] == nil {
				i, _ := t.RestType.Partition(uValueType)
				if (VoidType{}).Subsumes(i) {
					return VoidType{}, t
				}
				propTypeMap[prop] = i
				allSubsumed = false
			}
		}
		i, c := t.RestType.Partition(u.RestType)
		allSubsumed = allSubsumed && (VoidType{}).Subsumes(c)
		var complement Type
		if allSubsumed {
			complement = VoidType{}
		} else {
			complement = t // TODO give a more fine-grained type,
			// e.g. partitioning Obj<a: Num|Str> with Obj<a: Num> should give Obj<a: Str> as complement
		}
		return ObjType{
			PropTypeMap: propTypeMap,
			RestType:    i,
		}, complement
	case UnionType:
		return u.inversePartition(t)
	case AnyType:
		return t, VoidType{}
	default:
		return VoidType{}, t
	}
}

func (t ObjType) Instantiate(bindings map[string]Type) Type {
	propTypeMap := make(map[string]Type)
	for prop, valueType := range t.PropTypeMap {
		propTypeMap[prop] = valueType.Instantiate(bindings)
	}
	return ObjType{
		PropTypeMap: propTypeMap,
		RestType:    t.RestType.Instantiate(bindings),
	}
}

func (t ObjType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("Obj<")
	props := make([]string, len(t.PropTypeMap))
	i := 0
	for prop := range t.PropTypeMap {
		props[i] = prop
		i++
	}
	sort.Strings(props)
	if len(props) != 0 {
		first := true
		for _, prop := range props {
			if first {
				first = false
			} else {
				buffer.WriteString(", ")
			}
			buffer.WriteString(prop)
			buffer.WriteString(": ")
			buffer.WriteString(t.PropTypeMap[prop].String())
		}
	}
	if !t.RestType.Subsumes(AnyType{}) {
		if len(t.PropTypeMap) != 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(t.RestType.String())
	}
	buffer.WriteString(">")
	return buffer.String()
}

func (t ObjType) ElementType() Type {
	panic(t.String() + " is not a sequence type")
}

func (t ObjType) TypeForProp(prop string) Type {
	typeForProp := t.PropTypeMap[prop]
	if typeForProp == nil {
		return t.RestType
	}
	return typeForProp
}
