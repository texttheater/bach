package states

type State2 struct {
	fun       func() State2
	val       Value
	err       error
	Stack     *VariableStack
	TypeStack *BindingStack
}

var InitialState2 = &State2{
	val: NullValue{},
}

func (s State2) WithFunction(fun func() State2) State2 {
	return State2{
		fun:       fun,
		Stack:     s.Stack,
		TypeStack: s.TypeStack,
	}
}

func (s State2) WithValue(val Value) State2 {
	return State2{
		val:       val,
		Stack:     s.Stack,
		TypeStack: s.TypeStack,
	}
}

func (s State2) WithErrro(err error) State2 {
	return State2{
		err:       err,
		Stack:     s.Stack,
		TypeStack: s.TypeStack,
	}
}

func (s *State2) Eval() (Value, error) {
	for s.fun != nil {
		state := s.fun()
		s.fun = state.fun
		s.val = state.val
		s.err = state.err
	}
	return s.val, s.err
}

func (s *State2) EvalBool() (bool, error) {
	val, err := s.Eval()
	if err != nil {
		return false, err
	}
	return bool(val.(BoolValue)), nil
}

func (s *State2) EvalNum() (float64, error) {
	val, err := s.Eval()
	if err != nil {
		return 0, err
	}
	return float64(val.(NumValue)), nil
}

func (s *State2) EvalInt() (int, error) {
	val, err := s.Eval()
	if err != nil {
		return 0, err
	}
	return int(val.(NumValue)), nil
}

func (s *State2) EvalStr() (string, error) {
	val, err := s.Eval()
	if err != nil {
		return "", err
	}
	return string(val.(StrValue)), nil
}

func (s *State2) EvalArr() (*ArrValue, error) {
	val, err := s.Eval()
	if err != nil {
		return nil, err
	}
	return val.(*ArrValue), nil
}

func (s *State2) EvalObj() (ObjValue, error) {
	val, err := s.Eval()
	if err != nil {
		return nil, err
	}
	return val.(ObjValue), nil
}
