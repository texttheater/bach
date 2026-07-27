package states

type Action func(inputState State, args []Action) State

func SimpleAction(value Value) Action {
	return func(inputState State, args []Action) State {
		return State{
			Thunk:     ThunkFromValue(value),
			Stack:     inputState.Stack,
			TypeStack: inputState.TypeStack,
		}
	}
}
