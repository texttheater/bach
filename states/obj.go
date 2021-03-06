package states

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/texttheater/bach/types"
)

var lid *regexp.Regexp = regexp.MustCompile(`^[\p{L}_][\p{L}_0-9]*$`)

func ObjValueFromMap(m map[string]Value) ObjValue {
	propThunkMap := make(map[string]*Thunk)
	for k, v := range m {
		propThunkMap[k] = ThunkFromValue(v)
	}
	return propThunkMap
}

type ObjValue map[string]*Thunk

func (v ObjValue) Repr() (string, error) {
	buffer := bytes.Buffer{}
	buffer.WriteString("{")
	firstWritten := false
	for k, w := range v {
		if firstWritten {
			buffer.WriteString(", ")
		}
		if lid.MatchString(k) {
			buffer.WriteString(k)
		} else {
			buffer.WriteString(fmt.Sprintf("%q", k))
		}
		buffer.WriteString(": ")
		res := w.Eval()
		if res.Error != nil {
			return "", res.Error
		}
		wString, err := res.Value.Repr()
		if err != nil {
			return "", err
		}
		buffer.WriteString(wString)
		firstWritten = true
	}
	buffer.WriteString("}")
	return buffer.String(), nil
}

func (v ObjValue) Str() (string, error) {
	return v.Repr()
}

func (v ObjValue) Inhabits(t types.Type, stack *BindingStack) (bool, error) {
	switch t := t.(type) {
	case types.ObjType:
		for prop := range t.PropTypeMap {
			if v[prop] == nil {
				return false, nil
			}
		}
		for gotProp, thunk := range v {
			res := thunk.Eval()
			if res.Error != nil {
				return false, res.Error
			}
			if ok, err := res.Value.Inhabits(t.TypeForProp(gotProp), stack); !ok {
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
			vRes := vThunk.Eval()
			if vRes.Error != nil {
				return false, vRes.Error
			}
			wRes := wThunk.Eval()
			if wRes.Error != nil {
				return false, wRes.Error
			}
			equal, err := vRes.Value.Equal(wRes.Value)
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
