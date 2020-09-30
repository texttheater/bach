package states

import (
	"github.com/texttheater/bach/types"
)

type NullValue struct {
}

func (v NullValue) Repr() (string, error) {
	return "null", nil
}

func (v NullValue) Out() (string, error) {
	return v.Repr()
}

func (v NullValue) Inhabits(t types.Type, stack *BindingStack) (bool, error) {
	switch t := t.(type) {
	case types.NullType:
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

func (v NullValue) Equal(w Value) (bool, error) {
	switch w.(type) {
	case NullValue:
		return true, nil
	default:
		return false, nil
	}
}
