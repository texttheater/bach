package builtin

import (
	"fmt"
	"os"

	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initIO() {
	InitialShape.FuncerStack = InitialShape.FuncerStack.PushAll([]shapes.Funcer{
		shapes.SimpleFuncer(
			types.AnyType,
			"in",
			nil,
			types.ReaderType,
			func(inputValue values.Value, argValues []values.Value) values.Value {
				return values.ReaderValue{
					os.Stdin,
				}
			}),
		func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*shapes.Parameter, types.Type, states.Action, bool) {
			if gotName != "out" {
				return nil, nil, nil, false
			}
			if gotNumArgs != 0 {
				return nil, nil, nil, false
			}
			outputType := gotInputType
			action := func(inputState states.State, args []states.Action) states.State {
				fmt.Println(inputState.Value.Out())
				return inputState
			}
			return nil, outputType, action, true
		},
		func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*shapes.Parameter, types.Type, states.Action, bool) {
			if gotName != "err" {
				return nil, nil, nil, false
			}
			if gotNumArgs != 0 {
				return nil, nil, nil, false
			}
			outputType := gotInputType
			action := func(inputState states.State, args []states.Action) states.State {
				fmt.Fprintln(os.Stderr, inputState.Value)
				return inputState
			}
			return nil, outputType, action, true
		},
	})
}
