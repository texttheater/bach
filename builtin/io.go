package builtin

import (
	"bufio"
	"fmt"
	"os"

	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initIO() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.SimpleFuncer(
			types.AnyType{},
			"in",
			nil,
			types.ReaderType{},
			func(inputValue values.Value, argValues []values.Value) values.Value {
				return values.ReaderValue{
					os.Stdin,
				}
			},
		),
		functions.SimpleFuncer(
			types.ReaderType{},
			"lines",
			nil,
			&types.ArrType{types.StrType{}},
			func(inputValue values.Value, argValues []values.Value) values.Value {
				reader, _ := inputValue.(values.ReaderValue)
				scanner := bufio.NewScanner(reader.Reader)
				var next func() (values.Value, *values.ArrValue)
				next = func() (values.Value, *values.ArrValue) {
					ok := scanner.Scan()
					if !ok {
						return nil, nil
					}
					return values.StrValue(scanner.Text()), &values.ArrValue{
						Func: next,
					}
				}
				return &values.ArrValue{
					Func: next,
				}
			},
		),
		func(gotInputShape functions.Shape, gotCall functions.CallExpression, gotParams []*functions.Parameter) (functions.Shape, states.Action, bool, error) {
			if len(gotCall.Args)+len(gotParams) != 0 {
				return functions.Shape{}, nil, false, nil
			}
			if gotCall.Name != "out" {
				return functions.Shape{}, nil, false, nil
			}
			outputShape := functions.Shape{
				Type:  gotInputShape.Type,
				Stack: gotInputShape.Stack,
			}
			action := func(inputState states.State, args []states.Action) states.State {
				fmt.Println(inputState.Value.Out())
				return inputState
			}
			return outputShape, action, true, nil
		},
		func(gotInputShape functions.Shape, gotCall functions.CallExpression, gotParams []*functions.Parameter) (functions.Shape, states.Action, bool, error) {
			if len(gotCall.Args)+len(gotParams) != 0 {
				return functions.Shape{}, nil, false, nil
			}
			if gotCall.Name != "err" {
				return functions.Shape{}, nil, false, nil
			}
			outputShape := functions.Shape{
				Type:  gotInputShape.Type,
				Stack: gotInputShape.Stack,
			}
			action := func(inputState states.State, args []states.Action) states.State {
				fmt.Fprintln(os.Stderr, inputState.Value)
				return inputState
			}
			return outputShape, action, true, nil
		},
	})
}
