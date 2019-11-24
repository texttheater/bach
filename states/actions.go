package states

import (
	"github.com/texttheater/bach/values"
)

type Action func(inputState State, args []Action) Thunk

// TODO remove SetArg; currently one builtin funcer uses it; we should provide
// a better abstraction to builtin funcers making sure that partial application
// is handled correctly.

func (a Action) SetArg(arg Action) Action {
	return func(inputState State, outerArgs []Action) Thunk {
		args := make([]Action, len(outerArgs)+1)
		args[0] = arg
		for i, outerArg := range outerArgs {
			args[i+1] = outerArg
		}
		return a(inputState, args)
	}
}

func SimpleAction(value values.Value) Action {
	return func(inputState State, args []Action) Thunk {
		return EagerThunk(State{
			Value:     value,
			Stack:     inputState.Stack,
			TypeStack: inputState.TypeStack,
		}, false, nil)
	}
}
