package states

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type State struct {
	Error     error
	Drop      bool
	Value     values.Value
	Stack     *VariableStack
	TypeStack *BindingStack
}

type VariableStack struct {
	Head Variable
	Tail *VariableStack
}

func (s *VariableStack) Push(element Variable) *VariableStack {
	return &VariableStack{
		Head: element,
		Tail: s,
	}
}

type Variable struct {
	ID     interface{}
	Action Action
}

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

type Binding struct {
	Name string
	Type types.Type
}

var InitialState = State{
	Value: &values.NullValue{},
}
