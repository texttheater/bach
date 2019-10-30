package states

import (
	"github.com/texttheater/bach/values"
)

type State struct {
	Error     error
	Drop      bool
	Value     values.Value
	Stack     *VariableStack
	TypeStack *values.BindingStack
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

var InitialState = State{
	Value: values.NullValue{},
}
