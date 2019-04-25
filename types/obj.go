package types

import (
	"bytes"
	"sort"
)

type ObjType struct {
	Props       []string
	PropTypeMap map[string]Type
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
	if (VoidType{}).Subsumes(u) {
		return true
	}
	switch u := u.(type) {
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
		for _, disjunct := range u {
			if !t.Subsumes(disjunct) {
				return false
			}
		}
		return true
	default:
		return false
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
