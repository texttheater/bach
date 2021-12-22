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

func NumFromAction(state State, action Action) (float64, error) {
	res := action(state, nil).Eval()
	if res.Error != nil {
		return 0.0, res.Error
	}
	return float64(res.Value.(NumValue)), nil
}
