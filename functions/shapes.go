package functions

import (
	"github.com/texttheater/bach/types"
)

///////////////////////////////////////////////////////////////////////////////

type Shape struct {
	Type  types.Type
	Stack *FunctionStack
}

///////////////////////////////////////////////////////////////////////////////

type FunctionStack struct {
	Head *Function
	Tail *FunctionStack
}

func (stack *FunctionStack) Push(f *Function) *FunctionStack {
	return &FunctionStack{f, stack}
}

func (stack *FunctionStack) Pop() *FunctionStack {
	return stack.Tail
}

///////////////////////////////////////////////////////////////////////////////
