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
		functions.SimpleFuncer(
			types.StrType{},
			"split",
			nil,
			&types.ArrType{types.StrType{}},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				str, _ := inputValue.(states.StrValue)
				fields := strings.Fields(string(str))
				i := 0
				var next func() (states.Value, *states.ArrValue, error)
				next = func() (states.Value, *states.ArrValue, error) {
					if i >= len(fields) {
						return nil, nil, nil
					}
					head := states.StrValue(fields[i])
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
		functions.SimpleFuncer(
			&types.ArrType{types.StrType{}},
			"join",
			nil, // TODO separator
			types.StrType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				arr, _ := inputValue.(*states.ArrValue)
				err := arr.Eval()
				if err != nil {
					return nil, err
				}
				if arr.Head == nil {
					return states.StrValue(""), nil
				}
				var buffer bytes.Buffer
				str, err := arr.Head.Out()
				if err != nil {
					return nil, err
				}
				buffer.WriteString(str)
				arr = arr.Tail
				err = arr.Eval()
				if err != nil {
					return nil, err
				}
				for arr.Head != nil {
					// TODO separator
					str, err = arr.Head.Out()
					if err != nil {
						return nil, err
					}
					buffer.WriteString(str)
					arr = arr.Tail
					err = arr.Eval()
					if err != nil {
						return nil, err
					}
				}
				return states.StrValue(buffer.String()), nil
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
	})
}
