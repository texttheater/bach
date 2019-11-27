package states

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

var InitialState = State{
	Value: NullValue{},
}
