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

func (t objType) Subsumes(u Type) bool {
	if VoidType.Subsumes(u) {
		return true
	}
	switch u := u.(type) {
	case objType:
		for k, v1 := range t.propTypeMap {
			v2, ok := u.propTypeMap[k]
			if !ok {
				return false
			}
			if !v1.Subsumes(v2) {
				return false
			}
		}
		return true
	case unionType:
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
