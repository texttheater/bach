package functions

import (
	"github.com/texttheater/bach/states"
)

type Action func(inputStates states.State, outerArgs []Action) states.State

func (a Action) SetArg(arg Action) Action {
	return func(inputState states.State, outerArgs []Action) states.State {
		args := make([]Action, 0, len(outerArgs)+1)
		args = append(args, arg)
		for _, outerArg := range outerArgs {
			args = append(args, outerArg)
		}
		return a(inputState, args)
	}
}
