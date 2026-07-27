package states

type Action func(inputState State, args []Action) State

func SimpleAction(thunk *Thunk) Action {
	return func(inputState State, args []Action) State {
		return inputState.Replace(thunk)
	}
}
