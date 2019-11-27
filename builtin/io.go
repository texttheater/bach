package builtin

import (
	"bufio"
	"fmt"
	"os"

	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initIO() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.SimpleFuncer(
			types.AnyType{},
			"in",
			nil,
			types.ReaderType{},
			func(inputValue states.Value, argValues []states.Value) (states.Value, error) {
				return states.ReaderValue{
					os.Stdin,
				}, nil
			},
		),
		functions.SimpleFuncer(
			types.ReaderType{},
			"lines",
			nil,
			&types.ArrType{types.StrType{}},
			func(inputValue states.Value, argValues []states.Value) (states.Value, error) {
				reader, _ := inputValue.(states.ReaderValue)
				scanner := bufio.NewScanner(reader.Reader)
				var next func() (states.Value, *states.ArrValue, error)
				next = func() (states.Value, *states.ArrValue, error) {
					ok := scanner.Scan()
					if !ok {
						return nil, nil, nil
					}
					return states.StrValue(scanner.Text()), &states.ArrValue{
						Func: next,
					}, nil
				}
				return &states.ArrValue{
					Func: next,
				}, nil
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
			action := func(inputState states.State, args []states.Action) states.Thunk {
				str, err := inputState.Value.Out()
				if err != nil {
					return states.Thunk{State: states.State{}, Drop: false, Err: err}
				}
				fmt.Println(str)
				return states.Thunk{State: inputState, Drop: false, Err: nil}
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
			action := func(inputState states.State, args []states.Action) states.Thunk {
				str, err := inputState.Value.Out()
				if err != nil {
					return states.Thunk{State: states.State{}, Drop: false, Err: err}
				}
				fmt.Fprintln(os.Stderr, str)
				return states.Thunk{State: inputState, Drop: false, Err: nil}
			}
			return outputShape, action, true, nil
		},
	})
}
