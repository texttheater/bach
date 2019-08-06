package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initConv() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		func(gotInputShape functions.Shape, gotCall functions.CallExpression, gotParams []*functions.Parameter) (functions.Shape, states.Action, bool, error) {
			if len(gotCall.Args)+len(gotParams) != 0 {
				return functions.Shape{}, nil, false, nil
			}
			if gotCall.Name != "arr" {
				return functions.Shape{}, nil, false, nil
			}
			if !types.AnySeqType.Subsumes(gotInputShape.Type) {
				return functions.Shape{}, nil, false, nil
			}
			outputShape := functions.Shape{
				Type:  &types.ArrType{gotInputShape.Type.ElementType()},
				Stack: gotInputShape.Stack,
			}
			action := func(inputState states.State, args []states.Action) states.State {
				array := make([]values.Value, 0)
				for el := range inputState.Value.Iter() {
					array = append(array, el)
				}
				outputState := states.State{
					Value: values.ArrValue(array),
					Stack: inputState.Stack,
				}
				return outputState
			}
			return outputShape, action, true, nil
		},
	})
}
