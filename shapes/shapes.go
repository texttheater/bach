package shapes

import (
	"github.com/texttheater/bach/types"
)

type Shape struct {
	Type  types.Type
	Stack *FuncerStack
}

type FuncerStack struct {
	Head Funcer
	Tail *FuncerStack
}

func (s *FuncerStack) Push(funcer Funcer) *FuncerStack {
	return &FuncerStack{funcer, s}
}

func (s *FuncerStack) PushAll(funcers []Funcer) *FuncerStack {
	for i := len(funcers) - 1; i >= 0; i-- {
		s = s.Push(funcers[i])
	}
	return s
}
