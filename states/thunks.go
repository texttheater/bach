package states

type Thunk func() (State, bool, error, Thunk)

func EagerThunk(state State, drop bool, err error) Thunk {
	return func() (State, bool, error, Thunk) {
		return state, drop, err, nil
	}
}
