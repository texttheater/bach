package states

import (
	"github.com/texttheater/bach/types"
)

type Value interface {
	Repr() (string, error)
	Str() (string, error)
	Inhabits(types.Type, *BindingStack) (bool, error)
	Equal(Value) (bool, error)
}

func inhabits(v Value, t types.Union, stack *BindingStack) (bool, error) {
	for _, disjunct := range t {
		ok, err := v.Inhabits(disjunct, stack)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

func ValueFromAction(state State, action Action) (Value, error) {
	res := action(state, nil).Eval()
	if res.Error != nil {
		return nil, res.Error
	}
	return res.Value, nil
}
