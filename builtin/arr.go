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
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				arr, _ := inputValue.(*values.ArrValue)
				length := 0
				for arr != nil {
					length += 1
					arr = arr.Tail
				}
				return values.NumValue(length)
			},
		),
	})
}
