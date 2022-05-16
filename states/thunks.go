package states

type Thunk struct {
	Func      func() *Thunk
	Value     Value
	Error     error
	Stack     *VariableStack
	TypeStack *BindingStack
}

func (t *Thunk) Eval() (Value, error) {
	for t.Func != nil {
		thunk := t.Func()
		t.Func = thunk.Func
		t.Value = thunk.Value
		t.Error = thunk.Error
	}
	return t.Value, t.Error
}

func (t *Thunk) EvalNum() (float64, error) {
	val, err := t.Eval()
	if err != nil {
		return 0, err
	}
	return float64(val.(NumValue)), nil
}

func (t *Thunk) EvalInt() (int, error) {
	val, err := t.Eval()
	if err != nil {
		return 0, err
	}
	return int(val.(NumValue)), nil
}

func (t *Thunk) EvalStr() (string, error) {
	val, err := t.Eval()
	if err != nil {
		return "", err
	}
	return string(val.(StrValue)), nil
}

func (t *Thunk) EvalArr() (*ArrValue, error) {
	val, err := t.Eval()
	if err != nil {
		return nil, err
	}
	return val.(*ArrValue), nil
}

func ThunkFromValue(v Value) *Thunk {
	return ThunkFromState(State{
		Value: v,
	})
}

func ThunkFromError(err error) *Thunk {
	return &Thunk{
		Error: err,
	}
}

func ThunkFromState(state State) *Thunk {
	return &Thunk{
		Value:     state.Value,
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
	val, err := action(state, nil).Eval()
	if err != nil {
		return IterFromError(err)
	}
	return IterFromValue(val)
}
