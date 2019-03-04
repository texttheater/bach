package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initTypes() {
	InitialShape.FuncerStack = InitialShape.FuncerStack.PushAll([]functions.Funcer{
		func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*functions.Parameter, types.Type, functions.Action, bool) {
			if gotName != "type" {
				return nil, nil, nil, false
			}
			if gotNumArgs != 0 {
				return nil, nil, nil, false
			}
			outputType := types.StrType
			action := func(inputState functions.State, args []functions.Action) functions.State {
				outputValue := values.StrValue(gotInputType.String())
				outputState := functions.State{
					Value: outputValue,
					Stack: inputState.Stack,
				}
				return outputState
			}
			return nil, outputType, action, true
		},
	})
}
