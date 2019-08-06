package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initSeq() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.SimpleFuncer(
			types.AnyType{},
			"range",
			[]types.Type{types.NumType{}, types.NumType{}},
			&types.SeqType{types.NumType{}},
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				start := argumentValues[0].(values.NumValue)
				end := argumentValues[1].(values.NumValue)
				channel := make(chan values.Value)
				go func() {
					for i := start; i < end; i++ {
						channel <- values.NumValue(i)
					}
					close(channel)
				}()
				return values.SeqValue{types.NumType{}, channel}
			},
		),
	})
}
