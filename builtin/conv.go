package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initConv() {
	InitialShape.FuncerStack = InitialShape.FuncerStack.PushAll([]functions.Funcer{
		func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*functions.Parameter, types.Type, functions.Action, bool) {
			if !types.AnySeqType.Subsumes(gotInputType) {
				return nil, nil, nil, false
			}
			if gotName != "arr" {
				return nil, nil, nil, false
			}
			if gotNumArgs != 0 {
				return nil, nil, nil, false
			}
			outputType := &types.ArrType{gotInputType.ElementType()}
			action := func(inputState functions.State, args []functions.Action) functions.State {
				array := make([]values.Value, 0)
				for el := range inputState.Value.Iter() {
					array = append(array, el)
				}
				outputState := functions.State{
					Value: values.ArrValue(array),
					Stack: inputState.Stack,
				}
				return outputState
			}
			return nil, outputType, action, true
		},
		func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*functions.Parameter, types.Type, functions.Action, bool) {
			if gotName != "id" {
				return nil, nil, nil, false
			}
			if gotNumArgs != 0 {
				return nil, nil, nil, false
			}
			outputType := gotInputType
			action := func(inputState functions.State, args []functions.Action) functions.State {
				return inputState
			}
			return nil, outputType, action, true
		},
	})
}
