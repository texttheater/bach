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
	res := a(inputState, args).Eval()
	if res.Error != nil {
		return nil, res.Error
	}
	return res.Value, nil
}
