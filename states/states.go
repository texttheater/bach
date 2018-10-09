package states

import (
	"github.com/texttheater/bach/values"
)

type State struct {
	Value values.Value
	Stack *Stack
}

type Stack struct {
	Head NamedValue
	Tail *Stack
}

func (stack *Stack) Push(n NamedValue) *Stack {
	return &Stack{n, stack}
}

func (stack *Stack) Pop() *Stack {
	return stack.Tail
}

type NamedValue struct {
	Name  string
	Value values.Value
}

var InitialState = State{&values.NullValue{}, nil}
