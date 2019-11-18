package values

import (
	"github.com/texttheater/bach/types"
)

type BindingStack struct {
	Head Binding
	Tail *BindingStack
}

func (s *BindingStack) Push(element Binding) *BindingStack {
	return &BindingStack{
		Head: element,
		Tail: s,
	}
}

func (s *BindingStack) Inhabits(v Value, t types.TypeVariable) (bool, error) {
	if s == nil {
		return false, nil
	}
	if s.Head.Name == t.Name {
		return v.Inhabits(s.Head.Type, s)
	}
	return s.Tail.Inhabits(v, t)
}

type Binding struct {
	Name string
	Type types.Type
}
