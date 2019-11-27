package states

import (
	"bytes"

	"github.com/texttheater/bach/types"
)

func ObjValueFromMap(m map[string]Value) ObjValue {
	propThunkMap := make(map[string]*Thunk)
	for k, v := range m {
		propThunkMap[k] = ThunkFromValue(v)
	}
	return propThunkMap
}

type ObjValue map[string]*Thunk

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
		state, _, err := w.Eval()
		if err != nil {
			return "", err
		}
		wString, err := state.Value.String()
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
			thunk, ok := v[prop]
			if !ok {
				return false, nil
			}
			state, _, err := thunk.Eval()
			if err != nil {
				return false, err
			}
			if ok, err := state.Value.Inhabits(wantType, stack); !ok {
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
		for prop, vThunk := range v {
			wThunk, ok := w[prop]
			if !ok {
				return false, nil
			}
			vState, _, err := vThunk.Eval()
			if err != nil {
				return false, err
			}
			wState, _, err := wThunk.Eval()
			if err != nil {
				return false, err
			}
			equal, err := vState.Value.Equal(wState.Value)
			if err != nil {
				return false, err
			}
			if !equal {
				return false, nil
			}
		}
		return true, nil
	default:
		return false, nil
	}
}
