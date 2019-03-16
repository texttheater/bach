package functions

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
