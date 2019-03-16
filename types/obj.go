package types

import (
	"bytes"
	"sort"
)

func ObjType(propTypeMap map[string]Type) Type {
	props := make([]string, len(propTypeMap))
	i := 0
	for k := range propTypeMap {
		props[i] = k
		i++
	}
	sort.Strings(props)
	return objType{
		props:       props,
		propTypeMap: propTypeMap,
	}
}

type objType struct {
	props       []string
	propTypeMap map[string]Type
}

func (t objType) Subsumes(other Type) bool {
	if VoidType.Subsumes(other) {
		return true
	}
	otherObjType, ok := other.(objType)
	if !ok {
		return false
	}
	for k, v1 := range t.propTypeMap {
		v2, ok := otherObjType.propTypeMap[k]
		if !ok {
			return false
		}
		if !v1.Subsumes(v2) {
			return false
		}
	}
	return true
}

func (t objType) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString("Obj<")
	if len(t.props) > 0 {
		buffer.WriteString(t.props[0])
		buffer.WriteString(": ")
		buffer.WriteString(t.propTypeMap[t.props[0]].String())
		for _, prop := range t.props[1:] {
			typ := t.propTypeMap[prop]
			buffer.WriteString(", ")
			buffer.WriteString(prop)
			buffer.WriteString(": ")
			buffer.WriteString(typ.String())
		}
	}
	buffer.WriteString(">")
	return buffer.String()
}

func (t objType) ElementType() Type {
	panic(t.String() + " is not a sequence type")
}
