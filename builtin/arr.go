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
			"length",
			nil,
			types.NumType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				arr, _ := inputValue.(*values.ArrValue)
				length := arr.Length()
				return values.NumValue(length), nil
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
				var next func() (values.Value, *values.ArrValue)
				next = func() (values.Value, *values.ArrValue) {
					if i >= end {
						return nil, nil
					}
					head := values.NumValue(i)
					i++
					return head, &values.ArrValue{
						Func: next,
					}
				}
				return &values.ArrValue{
					Func: next,
				}, nil
			},
		),
	})
}
