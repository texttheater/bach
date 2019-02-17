package builtin

import (
	"fmt"
	"os"

	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

func initIO() {
	InitialShape.FuncerStack = InitialShape.FuncerStack.PushAll([]functions.Funcer{
		func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*functions.Parameter, types.Type, functions.Action, bool) {
			if gotName != "out" {
				return nil, nil, nil, false
			}
			if gotNumArgs != 0 {
				return nil, nil, nil, false
			}
			outputType := gotInputType
			action := func(inputState functions.State, args []functions.Action) functions.State {
				fmt.Println(inputState.Value)
				return inputState
			}
			return nil, outputType, action, true
		},
		func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*functions.Parameter, types.Type, functions.Action, bool) {
			if gotName != "err" {
				return nil, nil, nil, false
			}
			if gotNumArgs != 0 {
				return nil, nil, nil, false
			}
			outputType := gotInputType
			action := func(inputState functions.State, args []functions.Action) functions.State {
				fmt.Fprintln(os.Stderr, inputState.Value)
				return inputState
			}
			return nil, outputType, action, true
		},
	})
}
