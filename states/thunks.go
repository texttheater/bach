package states

type Thunk struct {
	Func  func() *Thunk
	State State
	Drop  bool
	Err   error
}

func (t *Thunk) Eval() (State, bool, error) {
	for t.Func != nil {
		thunk := t.Func()
		t.State = thunk.State
		t.Drop = thunk.Drop
		t.Err = thunk.Err
		t.Func = thunk.Func
	}
	return t.State, t.Drop, t.Err
}

func ThunkFromValue(v Value) *Thunk {
	return &Thunk{
		State: State{
			Value: v,
		},
	}
}
