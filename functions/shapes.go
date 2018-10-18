package functions

import (
	"github.com/texttheater/bach/types"
)

///////////////////////////////////////////////////////////////////////////////

type Shape struct {
	Type  types.Type
	Stack *NFFStack
}

///////////////////////////////////////////////////////////////////////////////

type NFFStack struct {
	Head NFF
	Tail *NFFStack
}

func (stack *NFFStack) Push(n NFF) *NFFStack {
	return &NFFStack{n, stack}
}

func (stack *NFFStack) Pop() *NFFStack {
	return stack.Tail
}

///////////////////////////////////////////////////////////////////////////////

type NFF struct {
	InputType types.Type // TODO type parameters
	Name      string     // TODO namespaces
	ArgTypes  []types.Type
	Funcer    func([]Function) Function // TODO first-class functions
}

///////////////////////////////////////////////////////////////////////////////
