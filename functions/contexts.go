package functions

import (
	"github.com/texttheater/bach/types"
)

type Context struct {
	Type          types.Type
	FunctionStack *FunctionStack
}

type FunctionStack struct {
	Head Function
	Tail *FunctionStack
}

func (s *FunctionStack) Push(function Function) *FunctionStack {
	return &FunctionStack{function, s}
}
