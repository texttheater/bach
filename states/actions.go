package states

import (
	"github.com/texttheater/bach/values"
)

type Action func(inputState State, args []Action) State

func (a Action) SetArg(arg Action) Action {
	return func(inputState State, outerArgs []Action) State {
		args := make([]Action, len(outerArgs)+1)
		args[0] = arg
		for i, outerArg := range outerArgs {
			args[i+1] = outerArg
		}
		return a(inputState, args)
	}
}

func SimpleAction(value values.Value) Action {
	return func(inputState State, args []Action) State {
		return State{
			Value:     value,
			Stack:     inputState.Stack,
			TypeStack: inputState.TypeStack,
		}
	}
}
