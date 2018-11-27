package functions

import (
	"github.com/texttheater/bach/values"
)

type State struct {
	Value values.Value
	Stack *VariableStack
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
	Id     interface{}
	Action Action
}

var InitialState State = State{
	Value: &values.NullValue{},
	Stack: nil,
}
