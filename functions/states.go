package functions

import (
	"github.com/texttheater/bach/values"
)

type State struct {
	Value values.Value
	Stack *VarStack
}

type VarStack struct {
	Head NamedValue
	Tail *VarStack
}

func (stack *VarStack) Push(n NamedValue) *VarStack {
	return &VarStack{n, stack}
}

func (stack *VarStack) Pop() *VarStack {
	return stack.Tail
}

type NamedValue struct {
	Name  string
	Value values.Value
}

var InitialState = State{&values.NullValue{}, nil}
