package values

import (
	"github.com/texttheater/bach/types"
)

type NullValue struct {
}

func (v NullValue) String() string {
	return "null"
}

func (v NullValue) Out() string {
	return v.String()
}

func (v NullValue) Inhabits(t types.Type, stack *BindingStack) bool {
	switch t := t.(type) {
	case types.NullType:
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

func (v NullValue) Equal(w Value) bool {
	switch w.(type) {
	case NullValue:
		return true
	default:
		return false
	}
}
