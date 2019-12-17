package states

import (
	"fmt"
)

type State struct {
	Value     Value
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

type IDStack struct {
	Head interface{}
	Tail *IDStack
}

func (s *IDStack) Contains(element interface{}) bool {
	for s != nil {
		if s.Head == element {
			return true
		}
		s = s.Tail
	}
	return false
}

func (s *IDStack) Add(element interface{}) *IDStack {
	if s.Contains(element) {
		return s
	}
	return &IDStack{
		Head: element,
		Tail: s,
	}
}

func (s *IDStack) AddAll(t *IDStack) *IDStack {
	for t != nil {
		s = s.Add(t.Head)
		t = t.Tail
	}
	return s
}

func (s *IDStack) String() string {
	var slice []interface{}
	for s != nil {
		slice = append(slice, s.Head)
		s = s.Tail
	}
	return fmt.Sprintf("%s", slice)
}

var InitialState = State{
	Value: NullValue{},
}
