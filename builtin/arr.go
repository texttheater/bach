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
				arr, _ := inputValue.(*states.ArrValue)
				err := arr.Eval()
				if err != nil {
					return nil, err
				}
				for arr.Head != nil {
					length += 1
					arr = arr.Tail
					err = arr.Eval()
					if err != nil {
						return nil, err
					}
				}
				return states.NumValue(float64(length)), nil
			},
		),
	})
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.SimpleFuncer(
			types.AnyType{},
			"range",
			[]types.Type{types.NumType{}, types.NumType{}},
			&types.ArrType{types.NumType{}},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				start := argumentValues[0].(states.NumValue)
				end := argumentValues[1].(states.NumValue)
				i := start
				var next func() (states.Value, *states.ArrValue, error)
				next = func() (states.Value, *states.ArrValue, error) {
					if i >= end {
						return nil, nil, nil
					}
					head := states.NumValue(i)
					i++
					return head, &states.ArrValue{
						Func: next,
					}, nil
				}
				return &states.ArrValue{
					Func: next,
				}, nil
			},
		),
	})
}
