package states

type Action func(inputState State, args []Action) *Thunk

func SimpleAction(value Value) Action {
	return func(inputState State, args []Action) *Thunk {
		return ThunkFromState(State{
			Value:     value,
			Stack:     inputState.Stack,
			TypeStack: inputState.TypeStack,
		})
	}
}
