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
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				str, _ := inputValue.(values.StrValue)
				fields := strings.Fields(string(str))
				i := 0
				var next func() (values.Value, *values.ArrValue, error)
				next = func() (values.Value, *values.ArrValue, error) {
					if i >= len(fields) {
						return nil, nil, nil
					}
					head := values.StrValue(fields[i])
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
