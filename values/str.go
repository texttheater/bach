package values

import (
	"fmt"

	"github.com/texttheater/bach/types"
)

type StrValue string

func (v StrValue) String() string {
	return fmt.Sprintf("%q", string(v))
}

func (v StrValue) Out() string {
	return string(v)
}

func (v StrValue) Inhabits(t types.Type, stack *BindingStack) bool {
	switch t := t.(type) {
	case types.StrType:
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

func (v StrValue) Equal(w Value) bool {
	switch w := w.(type) {
	case StrValue:
		return v == w
	default:
		return false
	}
}
