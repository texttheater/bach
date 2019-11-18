package values

import (
	"fmt"

	"github.com/texttheater/bach/types"
)

type StrValue string

func (v StrValue) String() (string, error) {
	return fmt.Sprintf("%q", string(v)), nil
}

func (v StrValue) Out() (string, error) {
	return string(v), nil
}

func (v StrValue) Inhabits(t types.Type, stack *BindingStack) (bool, error) {
	switch t := t.(type) {
	case types.StrType:
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

func (v StrValue) Equal(w Value) (bool, error) {
	switch w := w.(type) {
	case StrValue:
		return v == w, nil
	default:
		return false, nil
	}
}
