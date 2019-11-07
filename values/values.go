package values

import (
	"github.com/texttheater/bach/types"
)

type Value interface {
	String() string
	Out() string
	Iter() <-chan Value
	Inhabits(types.Type, *BindingStack) bool
	Equal(Value) bool
}

func inhabits(v Value, t types.UnionType, stack *BindingStack) bool {
	for _, disjunct := range t {
		if v.Inhabits(disjunct, stack) {
			return true
		}
	}
	return false
}
