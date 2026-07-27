package states

type Thunk struct {
	Func  func() (Value, error)
	Value Value
	Error error
}

func (t *Thunk) Eval() (Value, error) {
	if t.Func != nil {
		t.Value, t.Error = t.Func()
		t.Func = nil
	}
	return t.Value, t.Error
}

func (t *Thunk) EvalBool() (bool, error) {
	val, err := t.Eval()
	if err != nil {
		return false, err
	}
	return bool(val.(BoolValue)), nil
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

func (t *Thunk) EvalObj() (ObjValue, error) {
	val, err := t.Eval()
	if err != nil {
		return nil, err
	}
	return val.(ObjValue), nil
}

func ThunkFromFunc(fun func() (Value, error)) *Thunk {
	return &Thunk{
		Func: fun,
	}
}

func ThunkFromValue(val Value) *Thunk {
	return &Thunk{
		Value: val,
	}
}

func ThunkFromError(err error) *Thunk {
	return &Thunk{
		Error: err,
	}
}

func IterFromError(err error) func() (*Thunk, bool, error) {
	return func() (*Thunk, bool, error) {
		return nil, false, err
	}
}

func IterFromAction(state State, action Action) func() (*Thunk, bool, error) {
	return IterFromThunk(action(state, nil).Thunk)
}

func StateFromError(err error) State {
	return State{
		Thunk: ThunkFromError(err),
	}
}
