package builtin

import (
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initConv() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]shapes.Funcer{
		func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*shapes.Parameter, types.Type, states.Action, bool) {
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
			return nil, outputType, action, true
		},
	})
}
