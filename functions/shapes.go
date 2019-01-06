package functions

import (
	"github.com/texttheater/bach/types"
)

type Shape struct {
	Type        types.Type
	FuncerStack *FuncerStack
}

type FuncerStack struct {
	Head Funcer
	Tail *FuncerStack
}

func (s *FuncerStack) Push(funcer Funcer) *FuncerStack {
	return &FuncerStack{funcer, s}
}
