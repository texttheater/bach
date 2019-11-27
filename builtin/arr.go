package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initArr() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.SimpleFuncer(
			types.AnyArrType,
			"len",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				length := 0
				arr := inputValue.(*states.ArrValue)
				for arr != nil {
					length++
					var err error
					arr, err = arr.GetTail()
					if err != nil {
						return nil, err
					}
				}
				return states.NumValue(float64(length)), nil
			},
		),
	})
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.RegularFuncer(
			types.AnyType{},
			"range",
			[]*functions.Parameter{
				functions.SimpleParam(types.NumType{}),
				functions.SimpleParam(types.NumType{}),
			},
			&types.ArrType{types.NumType{}},
			func(inputState states.State, args []states.Action) *states.Thunk {
				state0, _, err := args[0](inputState, nil).Eval()
				if err != nil {
					return &states.Thunk{
						Err: err,
					}
				}
				start := float64(state0.Value.(states.NumValue))
				state1, _, err := args[1](inputState, nil).Eval()
				if err != nil {
					return &states.Thunk{
						Err: err,
					}
				}
				end := float64(state1.Value.(states.NumValue))
				var next func(float64) *states.Thunk
				next = func(i float64) *states.Thunk {
					if i >= end {
						return states.ThunkFromValue((*states.ArrValue)(nil))
					}
					return states.ThunkFromValue(&states.ArrValue{
						Head: states.NumValue(i),
						Tail: &states.Thunk{
							Func: func() *states.Thunk {
								return next(i + 1)
							},
						},
					})

				}
				return next(start)
			},
		),
	})
}
