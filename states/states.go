package states

import (
	"github.com/texttheater/bach/values"
)

type State struct {
	Value       values.Value
	ActionStack *ActionStack
}

type ActionStack struct {
	Head Action
	Tail *ActionStack
}

func (s *ActionStack) Push(action Action) *ActionStack {
	return &ActionStack{
		Head: action,
		Tail: s,
	}
}

type Action struct {
	Name    string
	Execute func(inputState State, outerArgs []Action) State
}

func (a Action) SetArg(arg Action) Action {
	return Action{
		Name: a.Name,
		Execute: func(inputState State, outerArgs []Action) State {
			args := make([]Action, 0, len(outerArgs)+1)
			args = append(args, arg)
			for _, outerArg := range outerArgs {
				args = append(args, outerArg)
			}
			return a.Execute(inputState, args)
		}}
}
