package states

import (
	"fmt"

	"github.com/texttheater/bach/types"
)

type StrValue string

func (v StrValue) Repr() (string, error) {
	return fmt.Sprintf("%q", string(v)), nil
}

func (v StrValue) Str() (string, error) {
	return string(v), nil
}

func (v StrValue) Inhabits(t types.Type, stack *BindingStack) (bool, error) {
	switch t := t.(type) {
	case types.Str:
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

func (v StrValue) Equal(w Value) (bool, error) {
	switch w := w.(type) {
	case StrValue:
		return v == w, nil
	default:
		return false, nil
	}
}
