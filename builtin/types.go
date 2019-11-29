package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initTypes() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		func(gotInputShape functions.Shape, gotCall functions.CallExpression, gotParams []*functions.Parameter) (functions.Shape, states.Action, bool, error) {
			if len(gotCall.Args)+len(gotParams) != 0 {
				return functions.Shape{}, nil, false, nil
			}
			if gotCall.Name != "type" {
				return functions.Shape{}, nil, false, nil
			}
			outputShape := functions.Shape{
				Type:  types.StrType{},
				Stack: gotInputShape.Stack,
			}
			action := func(inputState states.State, args []states.Action) *states.Thunk {
				outputValue := states.StrValue(gotInputShape.Type.String())
				outputState := states.State{
					Value:     outputValue,
					Stack:     inputState.Stack,
					TypeStack: inputState.TypeStack,
				}
				return states.ThunkFromState(outputState)
			}
			return outputShape, action, true, nil
		},
	})
}
