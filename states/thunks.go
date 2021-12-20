package states

type Thunk struct {
	Func      func() *Thunk
	Result    Result
	Stack     *VariableStack
	TypeStack *BindingStack
}

type Result struct {
	Value Value
	Drop  bool
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

func IterFromValue(v Value) func() (Value, bool, error) {
	return func() (Value, bool, error) {
		arr := v.(*ArrValue)
		if arr == nil {
			return nil, false, nil
		}
		var err error
		v, err = arr.GetTail()
		if err != nil {
			return nil, false, err
		}
		return arr.Head, true, nil
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

func ThunkFromIter(iter func() (Value, bool, error)) *Thunk {
	value, ok, err := iter()
	if err != nil {
		return ThunkFromError(err)
	}
	if !ok {
		return ThunkFromValue((*ArrValue)(nil))
	}
	return ThunkFromValue(&ArrValue{
		Head: value,
		Tail: &Thunk{
			Func: func() *Thunk {
				return ThunkFromIter(iter)
			},
		},
	})
}
