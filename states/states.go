package states

import (
	"github.com/texttheater/bach/values"
)

type State struct {
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
	Action Action // TODO this should be just a Value, right?
}

var InitialState = State{
	Value: values.NullValue{},
}
