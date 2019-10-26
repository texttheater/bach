package values

import (
	"github.com/texttheater/bach/types"
)

type SeqValue struct {
	ElementType types.Type
	Channel     chan Value
}

func (v SeqValue) String() string {
	return "<seq>"
}

func (v SeqValue) Out() string {
	return v.String()
}

func (v SeqValue) Iter() <-chan Value {
	// TODO safeguard against iterating twice?
	return v.Channel
}

func (v SeqValue) Inhabits(t types.Type, stack *BindingStack) bool {
	switch t := t.(type) {
	case *types.SeqType:
		return t.ElType.Subsumes(v.ElementType)
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
