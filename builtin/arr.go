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
				iter := states.IterFromValue(inputValue)
				for {
					_, ok, err := iter()
					if err != nil {
						return nil, err
					}
					if !ok {
						return states.NumValue(float64(length)), nil
					}
					length += 1
				}
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
				res0 := args[0](inputState, nil).Eval()
				if res0.Error != nil {
					return states.ThunkFromError(res0.Error)
				}
				start := float64(res0.State.Value.(states.NumValue))
				res1 := args[1](inputState, nil).Eval()
				if res1.Error != nil {
					return states.ThunkFromError(res1.Error)
				}
				end := float64(res1.State.Value.(states.NumValue))
				i := start
				var next func() (states.Value, bool, error)
				next = func() (states.Value, bool, error) {
					if i > end {
						return nil, false, nil
					}
					v := states.NumValue(i)
					i++
					return v, true, nil
				}
				return states.ThunkFromIter(next)
			},
		),
	})
}
