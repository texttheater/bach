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

func (a Action) Eval(inputState State, args []Action) (Value, error) {
	val, err := a(inputState, args).Eval()
	if err != nil {
		return nil, err
	}
	return val, nil
}
