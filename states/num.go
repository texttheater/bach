package states

import (
	"strconv"

	"github.com/texttheater/bach/types"
)

type NumValue float64

func (v NumValue) Repr() (string, error) {
	return strconv.FormatFloat(float64(v), 'g', -1, 64), nil
}

func (v NumValue) Str() (string, error) {
	return v.Repr()
}

func (v NumValue) Inhabits(t types.Type, stack *BindingStack) (bool, error) {
	switch t := t.(type) {
	case types.NumType:
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

func (v NumValue) Equal(w Value) (bool, error) {
	switch w := w.(type) {
	case NumValue:
		return v == w, nil
	default:
		return false, nil
	}
}
