package functions

type Action func(inputState State, outerArgs []Action) State

func (a Action) SetArg(arg Action) Action {
	return func(inputState State, outerArgs []Action) State {
		args := make([]Action, 0, len(outerArgs)+1)
		args = append(args, arg)
		for _, outerArg := range outerArgs {
			args = append(args, outerArg)
		}
		return a(inputState, args)
	}
}
