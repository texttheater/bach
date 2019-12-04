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
				str := string(inputState.Value.(states.StrValue))
				fields := strings.Fields(str)
				var iter func() (states.Value, bool, error)
				i := 0
				iter = func() (states.Value, bool, error) {
					if i >= len(fields) {
						return nil, false, nil
					}
					v := states.StrValue(fields[i])
					i++
					return v, true, nil
				}
				return states.ThunkFromIter(iter)
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
			&types.ArrType{types.StrType{}},
			"join",
			[]types.Type{types.StrType{}},
			types.StrType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				iter := states.IterFromValue(inputValue)
				sep := string(argumentValues[0].(states.StrValue))
				var buffer bytes.Buffer
				firstWritten := false
				for {
					value, ok, err := iter()
					if err != nil {
						return nil, err
					}
					if !ok {
						return states.StrValue(buffer.String()), nil
					}
					if firstWritten {
						buffer.WriteString(sep)
					}
					buffer.WriteString(string(value.(states.StrValue)))
					firstWritten = true
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
