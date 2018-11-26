package states

import (
	"github.com/texttheater/bach/values"
)

type State struct {
	Value values.Value
	Stack *Stack
}

type Stack struct {
	Head Variable
	Tail *Stack
}

func (s *Stack) Push(element Variable) *Stack {
	return &Stack{
		Head: element,
		Tail: s,
	}
}

type Variable struct {
	Name  string
	Value values.Value
}

var InitialState State = State{
	Value: &values.NullValue{},
	Stack: nil,
}
