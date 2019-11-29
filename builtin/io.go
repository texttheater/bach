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
		functions.RegularFuncer(
			types.ReaderType{},
			"lines",
			nil,
			&types.ArrType{types.StrType{}},
			func(inputState states.State, args []states.Action) *states.Thunk {
				reader := inputState.Value.(states.ReaderValue)
				scanner := bufio.NewScanner(reader.Reader)
				output := make(chan states.Result)
				go func() {
					defer close(output)
					for {
						ok := scanner.Scan()
						if !ok {
							break
						}
						output <- states.Result{
							State: states.State{
								Value: states.StrValue(scanner.Text()),
							},
						}
					}
				}()
				return states.ThunkFromChannel(output)
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
			action := func(inputState states.State, args []states.Action) *states.Thunk {
				str, err := inputState.Value.Out()
				if err != nil {
					return states.ThunkFromError(err)

				}
				fmt.Println(str)
				return states.ThunkFromState(inputState)
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
			action := func(inputState states.State, args []states.Action) *states.Thunk {
				str, err := inputState.Value.Out()
				if err != nil {
					return states.ThunkFromError(err)

				}
				fmt.Fprintln(os.Stderr, str)
				return states.ThunkFromState(inputState)
			}
			return outputShape, action, true, nil
		},
	})
}
