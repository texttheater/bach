package states

import (
	"strconv"

	"github.com/texttheater/bach/types"
)

type BoolValue bool

func (v BoolValue) Repr() (string, error) {
	return strconv.FormatBool(bool(v)), nil
}

func (v BoolValue) Str() (string, error) {
	return v.Repr()
}

func (v BoolValue) Data() (any, error) {
	return v, nil
}

func (v BoolValue) Inhabits(t types.Type, stack *BindingStack) (bool, error) {
	switch t := t.(type) {
	case types.Bool:
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

func (v BoolValue) Equal(w Value) (bool, error) {
	switch w := w.(type) {
	case BoolValue:
		return v == w, nil
	default:
		return false, nil
	}
}
