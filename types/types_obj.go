package types

import (
	"bytes"
	"sort"
)

type ObjType struct {
	Props map[string]Type
	Rest  Type
}

var VoidObjType Type = ObjType{
	Props: make(map[string]Type),
	Rest:  VoidType{},
}

var AnyObjType Type = ObjType{
	Props: make(map[string]Type),
	Rest:  AnyType{},
}

func (t ObjType) Subsumes(u Type) bool {
	switch u := u.(type) {
	case VoidType:
		return true
	case ObjType:
		for prop, wantType := range t.Props {
			gotType, ok := u.Props[prop]
			if !ok {
				return false
			}
			if !wantType.Subsumes(gotType) {
				return false
			}
		}
		if !t.Rest.Subsumes(u.Rest) {
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
		for prop, wantType := range t.Props {
			gotType, ok := u.Props[prop]
			if !ok {
				return false
			}
			if !wantType.Bind(gotType, bindings) {
				return false
			}
		}
		restType := u.Rest
		for prop, gotType := range u.Props {
			if _, ok := t.Props[prop]; !ok {
				restType = NewUnionType(restType, gotType)
			}
		}
		if !t.Rest.Bind(restType, bindings) {
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
		for prop, tValueType := range t.Props {
			uValueType := u.Props[prop]
			if uValueType != nil {
				i, c := tValueType.Partition(uValueType)
				if (VoidType{}).Subsumes(i) {
					return VoidType{}, t
				}
				propTypeMap[prop] = i
				allSubsumed = allSubsumed && (VoidType{}).Subsumes(c)
			} else {
				i, c := tValueType.Partition(u.Rest)
				if (VoidType{}).Subsumes(i) {
					return VoidType{}, t
				}
				propTypeMap[prop] = i
				allSubsumed = allSubsumed && (VoidType{}).Subsumes(c)
			}
		}
		for prop, uValueType := range u.Props {
			if t.Props[prop] == nil {
				i, _ := t.Rest.Partition(uValueType)
				if (VoidType{}).Subsumes(i) {
					return VoidType{}, t
				}
				propTypeMap[prop] = i
				allSubsumed = false
			}
		}
		i, c := t.Rest.Partition(u.Rest)
		allSubsumed = allSubsumed && (VoidType{}).Subsumes(c)
		var complement Type
		if allSubsumed {
			complement = VoidType{}
		} else {
			complement = t // TODO give a more fine-grained type,
			// e.g. partitioning Obj<a: Num|Str> with Obj<a: Num> should give Obj<a: Str> as complement
		}
		return ObjType{
			Props: propTypeMap,
			Rest:  i,
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
	for prop, valueType := range t.Props {
		propTypeMap[prop] = valueType.Instantiate(bindings)
	}
	return ObjType{
		Props: propTypeMap,
		Rest:  t.Rest.Instantiate(bindings),
	}
}

func (t ObjType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("Obj<")
	props := make([]string, len(t.Props))
	i := 0
	for prop := range t.Props {
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
			buffer.WriteString(t.Props[prop].String())
		}
	}
	if len(t.Props) != 0 {
		buffer.WriteString(", ")
	}
	buffer.WriteString(t.Rest.String())
	buffer.WriteString(">")
	return buffer.String()
}

func (t ObjType) ElementType() Type {
	panic(t.String() + " is not a sequence type")
}

func (t ObjType) TypeForProp(prop string) Type {
	typeForProp := t.Props[prop]
	if typeForProp == nil {
		return t.Rest
	}
	return typeForProp
}
