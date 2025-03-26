package states

import (
	"fmt"
	"strings"

	"github.com/texttheater/bach/types"
)

type StrValue string

func (v StrValue) Repr() (string, error) {
	s := string(v)
	s = strings.Replace(s, "{", "{{", -1)
	s = strings.Replace(s, "}", "}}", -1)
	s = fmt.Sprintf("%q", s)
	return s, nil
}

func (v StrValue) Str() (string, error) {
	return string(v), nil
}

func (v StrValue) Data() (any, error) {
	return v, nil
}

func (v StrValue) Inhabits(t types.Type, stack *BindingStack) (bool, error) {
	switch t := t.(type) {
	case types.StrType:
		return true, nil
	case types.UnionType:
		return inhabits(v, t, stack)
	case types.AnyType:
		return true, nil
	case types.TypeVar:
		return stack.Inhabits(v, t)
	default:
		return false, nil
	}
}

func (v StrValue) Equal(w Value) (bool, error) {
	switch w := w.(type) {
	case StrValue:
		return v == w, nil
	default:
		return false, nil
	}
}
