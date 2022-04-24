package states

type Thunk struct {
	Func      func() *Thunk
	Result    Result
	Stack     *VariableStack
	TypeStack *BindingStack
}

// TODO make this type private, perhaps abolish it
type Result struct {
	Value Value
	Error error
}

func (t *Thunk) eval() Result {
	for t.Func != nil {
		thunk := t.Func()
		t.Func = thunk.Func
		t.Result = thunk.Result
	}
	return t.Result
}

// TODO check where this is called, see if we can call type-specific versions
// instead
func (t *Thunk) Eval() (Value, error) {
	res := t.eval()
	return res.Value, res.Error
}

func (t *Thunk) EvalNum() (float64, bool, error) {
	val, err := t.Eval()
	if err != nil {
		return 0, false, err
	}
	v, ok := val.(NumValue)
	if !ok {
		return 0, false, nil
	}
	return float64(v), true, nil
}

func (t *Thunk) EvalStr() (string, bool, error) {
	val, err := t.Eval()
	if err != nil {
		return "", false, err
	}
	v, ok := val.(StrValue)
	if !ok {
		return "", false, nil
	}
	return string(v), true, nil
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
	val, err := action(state, nil).Eval()
	if err != nil {
		return IterFromError(err)
	}
	return IterFromValue(val)
}
