package states

import (
	"bytes"

	"github.com/texttheater/bach/types"
)

type ObjValue map[string]Value

func (v ObjValue) String() (string, error) {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	firstWritten := false
	for k, w := range v {
		if firstWritten {
			buffer.WriteString(", ")
		}
		buffer.WriteString(k)
		buffer.WriteString(": ")
		wString, err := w.String()
		if err != nil {
			return "", err
		}
		buffer.WriteString(wString)
		firstWritten = true
	}
	buffer.WriteString("}")
	return buffer.String(), nil
}

func (v ObjValue) Out() (string, error) {
	return v.String()
}

func (v ObjValue) Inhabits(t types.Type, stack *BindingStack) (bool, error) {
	switch t := t.(type) {
	case types.ObjType:
		for prop, wantType := range t.PropTypeMap {
			gotValue, ok := v[prop]
			if !ok {
				return false, nil
			}
			if ok, err := gotValue.Inhabits(wantType, stack); !ok {
				return false, err
			}
		}
		return true, nil
	case types.UnionType:
		return inhabits(v, t, stack)
	case types.AnyType:
		return true, nil
	case types.TypeVariable:
		return stack.Inhabits(v, t)
	default:
		return false, nil
	}
}

func (v ObjValue) Equal(w Value) (bool, error) {
	switch w := w.(type) {
	case ObjValue:
		if len(v) != len(w) {
			return false, nil
		}
		for k, l := range v {
			m, ok := w[k]
			if !ok {
				return false, nil
			}
			if ok, err := l.Equal(m); !ok {
				return false, err
			}
		}
		return true, nil
	default:
		return false, nil
	}
}
