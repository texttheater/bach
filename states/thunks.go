package states

type Thunk struct {
	Func      func() *Thunk
	Result    Result
	Stack     *VariableStack
	TypeStack *BindingStack
}

type Result struct {
	Value Value
	Error error
}

func (t *Thunk) Eval() Result {
	for t.Func != nil {
		thunk := t.Func()
		t.Func = thunk.Func
		t.Result = thunk.Result
	}
	return t.Result
}

func ThunkFromValue(v Value) *Thunk {
	return ThunkFromState(State{
		Value: v,
	})
}

func ThunkFromError(err error) *Thunk {
	return &Thunk{
		Result: Result{
			Error: err,
		},
	}
}

func ThunkFromState(state State) *Thunk {
	return &Thunk{
		Result: Result{
			Value: state.Value,
		},
		Stack:     state.Stack,
		TypeStack: state.TypeStack,
	}
}

func IterFromError(err error) func() (Value, bool, error) {
	return func() (Value, bool, error) {
		return nil, false, err
	}
}

func IterFromAction(state State, action Action) func() (Value, bool, error) {
	res := action(state, nil).Eval()
	if res.Error != nil {
		return IterFromError(res.Error)
	}
	return IterFromValue(res.Value)
}
