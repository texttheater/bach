package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initArr() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.SimpleFuncer(
			types.AnyArrType,
			"len",
			nil,
			types.NumType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				length := 0
				arr, _ := inputValue.(*values.ArrValue)
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
				return values.NumValue(float64(length)), nil
			},
		),
	})
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.SimpleFuncer(
			types.AnyType{},
			"range",
			[]types.Type{types.NumType{}, types.NumType{}},
			&types.ArrType{types.NumType{}},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				start := argumentValues[0].(values.NumValue)
				end := argumentValues[1].(values.NumValue)
				i := start
				var next func() (values.Value, *values.ArrValue, error)
				next = func() (values.Value, *values.ArrValue, error) {
					if i >= end {
						return nil, nil, nil
					}
					head := values.NumValue(i)
					i++
					return head, &values.ArrValue{
						Func: next,
					}, nil
				}
				return &values.ArrValue{
					Func: next,
				}, nil
			},
		),
	})
}
