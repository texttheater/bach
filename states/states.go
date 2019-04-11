package states

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
	ID     interface{}
	Action Action
}

var InitialState = State{
	Value: values.NullValue(),
	Stack: nil,
}
