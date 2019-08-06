package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
)

func initValues() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		func(gotInputShape functions.Shape, gotCall functions.CallExpression, gotParams []*functions.Parameter) (functions.Shape, states.Action, bool, error) {
			if len(gotCall.Args)+len(gotParams) != 0 {
				return functions.Shape{}, nil, false, nil
			}
			if gotCall.Name != "id" {
				return functions.Shape{}, nil, false, nil
			}
			action := func(inputState states.State, args []states.Action) states.State {
				return inputState
			}
			return gotInputShape, action, true, nil
		},
	})
}
