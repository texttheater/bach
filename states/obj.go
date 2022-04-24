package states

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/texttheater/bach/types"
)

var lid *regexp.Regexp = regexp.MustCompile(`^[\p{L}_][\p{L}_0-9]*$`)
var op1 *regexp.Regexp = regexp.MustCompile(`[+\-*/%<>=]`)
var op2 *regexp.Regexp = regexp.MustCompile(`==|<=|>=|\*\*`)
var num *regexp.Regexp = regexp.MustCompile(`\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+`)

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
		if lid.MatchString(k) || op1.MatchString(k) || op2.MatchString(k) ||
			num.MatchString(k) {
			buffer.WriteString(k)
		} else {
			buffer.WriteString(fmt.Sprintf("%q", k))
		}
		buffer.WriteString(": ")
		val, err := w.Eval()
		if err != nil {
			return "", err
		}
		wString, err := val.Repr()
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
	case types.Obj:
		for prop := range t.Props {
			if v[prop] == nil {
				return false, nil
			}
		}
		for gotProp, thunk := range v {
			val, err := thunk.Eval()
			if err != nil {
				return false, err
			}
			if ok, err := val.Inhabits(t.TypeForProp(gotProp), stack); !ok {
				return false, err
			}
		}
		return true, nil
	case types.Union:
		return inhabits(v, t, stack)
	case types.Any:
		return true, nil
	case types.Var:
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
			vVal, vErr := vThunk.Eval()
			if vErr != nil {
				return false, vErr
			}
			wVal, wErr := wThunk.Eval()
			if wErr != nil {
				return false, wErr
			}
			equal, err := vVal.Equal(wVal)
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
