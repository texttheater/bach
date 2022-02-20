package states

import (
	"fmt"
)

type State struct {
	Value     Value
	Stack     *VariableStack
	TypeStack *BindingStack
}

// Clear returns an identical state except that the Value is replaced with null
// TODO do we still need this method?
func (s State) Clear() State {
	return State{
		Value:     NullValue{},
		Stack:     s.Stack,
		TypeStack: s.TypeStack,
	}
}

// Replace returns an identical state except that the Value is replaced with the given value
func (s State) Replace(v Value) State {
	return State{
		Value:     v,
		Stack:     s.Stack,
		TypeStack: s.TypeStack,
	}
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

func (s *VariableStack) Keep(ids *IDStack) *VariableStack {
	var first *VariableStack
	var current *VariableStack
	for s != nil {
		if ids.Contains(s.Head.ID) {
			next := &VariableStack{
				Head: s.Head,
			}
			if current == nil {
				first = next
				current = next
			} else {
				current.Tail = next
				current = next
			}
		}
		s = s.Tail
	}
	return first
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
