package shapes

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

type Shape struct {
	Type          types.Type
	FunctionStack *FunctionStack
}

type FunctionStack struct {
	Head functions.Function
	Tail *FunctionStack
}

func (s *FunctionStack) Push(function functions.Function) *FunctionStack {
	return &FunctionStack{function, s}
}
