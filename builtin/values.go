package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initValues() {
	InitialShape.FuncerStack = InitialShape.FuncerStack.PushAll([]functions.Funcer{
		func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*functions.Parameter, types.Type, states.Action, bool) {
			if gotName != "id" {
				return nil, nil, nil, false
			}
			if gotNumArgs != 0 {
				return nil, nil, nil, false
			}
			outputType := gotInputType
			action := func(inputState states.State, args []states.Action) states.State {
				return inputState
			}
			return nil, outputType, action, true
		},
	})
}
