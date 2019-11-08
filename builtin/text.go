package builtin

import (
	"strings"

	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initText() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.SimpleFuncer(
			types.StrType{},
			"split",
			nil,
			&types.ArrType{types.StrType{}},
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				str, _ := inputValue.(values.StrValue)
				fields := strings.Fields(string(str))
				channel := make(chan values.Value)
				go func() {
					for _, field := range fields {
						channel <- values.StrValue(field)
					}
					close(channel)
				}()
				return &values.ArrValue{
					Channel: channel,
				}
			},
		),
	})
}
