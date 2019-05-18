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
