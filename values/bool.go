package values

import (
	"strconv"

	"github.com/texttheater/bach/types"
)

type BoolValue bool

func (v BoolValue) String() (string, error) {
	return strconv.FormatBool(bool(v)), nil
}

func (v BoolValue) Out() (string, error) {
	return v.String()
}

func (v BoolValue) Inhabits(t types.Type, stack *BindingStack) (bool, error) {
	switch t := t.(type) {
	case types.BoolType:
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

func (v BoolValue) Equal(w Value) (bool, error) {
	switch w := w.(type) {
	case BoolValue:
		return v == w, nil
	default:
		return false, nil
	}
}
