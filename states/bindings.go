package states

import (
	"github.com/texttheater/bach/types"
)

type BindingStack struct {
	Head Binding
	Tail *BindingStack
}

func (s *BindingStack) Update(name string, t types.Type) *BindingStack {
	if s == nil {
		return &BindingStack{
			Head: Binding{
				Name: name,
				Type: t,
			},
		}
	}
	if s.Head.Name == name {
		return &BindingStack{
			Head: Binding{
				Name: name,
				Type: t,
			},
			Tail: s.Tail,
		}
	}
	return &BindingStack{
		Head: s.Head,
		Tail: s.Tail.Update(name, t),
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
