package states

type Thunk func() (State, bool, error, Thunk)

func (t Thunk) Eval() (State, bool, error) {
	var state State
	var drop bool
	var err error
	for t != nil {
		state, drop, err, t = t()
	}
	return state, drop, err
}

func EagerThunk(state State, drop bool, err error) Thunk {
	return func() (State, bool, error, Thunk) {
		return state, drop, err, nil
	}
}
