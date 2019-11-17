package values

import (
	"bytes"

	"github.com/texttheater/bach/types"
)

type ObjValue map[string]Value

func (v ObjValue) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	firstWritten := false
	for k, w := range v {
		if firstWritten {
			buffer.WriteString(", ")
		}
		buffer.WriteString(k)
		buffer.WriteString(": ")
		buffer.WriteString(w.String())
		firstWritten = true
	}
	buffer.WriteString("}")
	return buffer.String()
}

func (v ObjValue) Out() string {
	return v.String()
}

func (v ObjValue) Inhabits(t types.Type, stack *BindingStack) bool {
	switch t := t.(type) {
	case types.ObjType:
		for prop, wantType := range t.PropTypeMap {
			gotValue, ok := v[prop]
			if !ok {
				return false
			}
			if !gotValue.Inhabits(wantType, stack) {
				return false
			}
		}
		return true
	case types.UnionType:
		return inhabits(v, t, stack)
	case types.AnyType:
		return true
	case types.TypeVariable:
		return stack.Inhabits(v, t)
	default:
		return false
	}
}

func (v ObjValue) Equal(w Value) bool {
	switch w := w.(type) {
	case ObjValue:
		if len(v) != len(w) {
			return false
		}
		for k, l := range v {
			m, ok := w[k]
			if !ok {
				return false
			}
			if !l.Equal(m) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
