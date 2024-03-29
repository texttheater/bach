package states

import (
	"math"
	"strconv"

	"github.com/texttheater/bach/types"
)

type NumValue float64

func (v NumValue) Repr() (string, error) {
	f := float64(v)
	if math.IsInf(f, 1) {
		return "inf", nil
	}
	if math.IsInf(f, -1) {
		return "-inf", nil
	}
	if math.IsNaN(f) {
		return "nan", nil
	}
	return strconv.FormatFloat(f, 'g', -1, 64), nil
}

func (v NumValue) Str() (string, error) {
	return v.Repr()
}

func (v NumValue) Data() (any, error) {
	return v, nil
}

func (v NumValue) Inhabits(t types.Type, stack *BindingStack) (bool, error) {
	switch t := t.(type) {
	case types.Num:
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

func (v NumValue) Equal(w Value) (bool, error) {
	switch w := w.(type) {
	case NumValue:
		return v == w, nil
	default:
		return false, nil
	}
}
