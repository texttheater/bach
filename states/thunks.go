package states

type Thunk struct {
	Func   func() *Thunk
	Result Result
}

type Result struct {
	State State
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
			State: state,
		},
	}
}

func ChannelFromValue(value Value) <-chan Result {
	channel := make(chan Result)
	go func() {
		defer close(channel)
		arr := value.(*ArrValue)
		for arr != nil {
			channel <- Result{
				State: State{
					Value: arr.Head,
				},
			}
			var err error
			arr, err = arr.GetTail()
			if err != nil {
				channel <- Result{
					Error: err,
				}
			}
		}
	}()
	return channel
}

func ThunkFromChannel(channel <-chan Result) *Thunk {
	var next func() *Thunk
	next = func() *Thunk {
		res, ok := <-channel
		if !ok {
			return ThunkFromValue((*ArrValue)(nil))
		}
		if res.Error != nil {
			return ThunkFromError(res.Error)
		}
		return ThunkFromValue(&ArrValue{
			Head: res.State.Value,
			Tail: &Thunk{
				Func: func() *Thunk {
					return next()
				},
			},
		})
	}
	return next()
}
