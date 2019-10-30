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
			&types.SeqType{types.StrType{}},
			func(inputValue values.Value, argValues []values.Value) values.Value {
				reader, _ := inputValue.(values.ReaderValue)
				lines := make(chan values.Value)
				scanner := bufio.NewScanner(reader.Reader)
				go func() {
					for scanner.Scan() {
						lines <- values.StrValue(scanner.Text())
					}
					close(lines)
				}()
				return values.SeqValue{types.StrType{}, lines}
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
