package builtin

import (
	"bytes"
	"strings"

	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initText() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.RegularFuncer(
			types.StrType{},
			"split",
			nil,
			&types.ArrType{types.StrType{}},
			func(inputState states.State, args []states.Action) *states.Thunk {
				str, _ := inputState.Value.(states.StrValue)
				fields := strings.Fields(string(str))
				i := 0
				var next func() *states.Thunk
				next = func() *states.Thunk {
					if i >= len(fields) {
						return states.ThunkFromValue((*states.ArrValue)(nil))
					}
					return states.ThunkFromValue(
						&states.ArrValue{
							Head: states.StrValue(fields[i]),
							Tail: &states.Thunk{
								Func: func() *states.Thunk {
									i++
									return next()
								},
							},
						},
					)
				}
				return next()
			},
		),
		functions.SimpleFuncer(
			&types.ArrType{types.StrType{}},
			"join",
			nil,
			types.StrType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				iter := states.IterFromValue(inputValue)
				var buffer bytes.Buffer
				for {
					value, ok, err := iter()
					if err != nil {
						return nil, err
					}
					if !ok {
						return states.StrValue(buffer.String()), nil
					}
					buffer.WriteString(string(value.(states.StrValue)))
				}
			},
		),
		functions.SimpleFuncer(
			types.StrType{},
			"==",
			[]types.Type{types.StrType{}},
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				str1 := string(inputValue.(states.StrValue))
				str2 := string(argumentValues[0].(states.StrValue))
				return states.BoolValue(str1 == str2), nil
			},
		),
		functions.SimpleFuncer(
			types.StrType{},
			"+",
			[]types.Type{types.StrType{}},
			types.StrType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				str1 := string(inputValue.(states.StrValue))
				str2 := string(argumentValues[0].(states.StrValue))
				return states.StrValue(str1 + str2), nil
			},
		),
	})
}
