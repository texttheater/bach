package builtin

import (
	"bytes"
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
		functions.SimpleFuncer(
			&types.ArrType{types.StrType{}},
			"join",
			nil, // TODO separator
			types.StrType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				arr, _ := inputValue.(*values.ArrValue)
				err := arr.Eval()
				if err != nil {
					return nil, err
				}
				if arr.Head == nil {
					return values.StrValue(""), nil
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
				return values.StrValue(buffer.String()), nil
			},
		),
		functions.SimpleFuncer(
			types.StrType{},
			"==",
			[]types.Type{types.StrType{}},
			types.BoolType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				str1 := string(inputValue.(values.StrValue))
				str2 := string(argumentValues[0].(values.StrValue))
				return values.BoolValue(str1 == str2), nil
			},
		),
	})
}
