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
	propThunkMap := make(map[string]Value)
	for k, v := range m {
		propThunkMap[k] = v
	}
	return propThunkMap
}

type ObjValue map[string]Value

func (v ObjValue) Repr() (string, error) {
	buffer := bytes.Buffer{}
	buffer.WriteString("{")
	firstWritten := false
	for k, val := range v {
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

func (v ObjValue) Data() (any, error) {
	res := make(map[string]any, 0)
	for k, val := range v {
		data, err := val.Data()
		if err != nil {
			return nil, err
		}
		res[k] = data
	}
	return res, nil
}

func (v ObjValue) Inhabits(t types.Type, stack *BindingStack) (bool, error) {
	switch t := t.(type) {
	case types.Obj:
		for prop := range t.Props {
			if v[prop] == nil {
				return false, nil
			}
		}
		for gotProp, val := range v {
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
		for prop, vVal := range v {
			wVal, ok := w[prop]
			if !ok {
				return false, nil
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
