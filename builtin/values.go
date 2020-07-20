package builtin

import (
	"strconv"

	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initValues() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		func(gotInputShape functions.Shape, gotCall functions.CallExpression, gotParams []*parameters.Parameter) (functions.Shape, states.Action, *states.IDStack, bool, error) {
			if len(gotCall.Args)+len(gotParams) != 0 {
				return functions.Shape{}, nil, nil, false, nil
			}
			if gotCall.Name != "id" {
				return functions.Shape{}, nil, nil, false, nil
			}
			action := func(inputState states.State, args []states.Action) *states.Thunk {
				return states.ThunkFromState(inputState)
			}
			return gotInputShape, action, nil, true, nil
		},
		functions.SimpleFuncer(
			types.StrType{},
			"parseFloat",
			nil,
			types.NumType{},
			func(inputValue states.Value, argValues []states.Value) (states.Value, error) {
				s := string(inputValue.(states.StrValue))
				n, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return nil, err
				}
				return states.NumValue(n), nil
			},
		),
		functions.SimpleFuncer(
			types.StrType{},
			"parseInt",
			[]types.Type{
				types.NumType{},
			},
			types.NumType{},
			func(inputValue states.Value, argValues []states.Value) (states.Value, error) {
				s := string(inputValue.(states.StrValue))
				b := argValues[0].(states.NumValue)
				n, err := strconv.ParseInt(s, int(b), 64)
				if err != nil {
					return nil, err
				}
				return states.NumValue(n), nil
			},
		),
		functions.SimpleFuncer(
			types.StrType{},
			"parseInt",
			nil,
			types.NumType{},
			func(inputValue states.Value, argValues []states.Value) (states.Value, error) {
				s := string(inputValue.(states.StrValue))
				b := 10
				n, err := strconv.ParseInt(s, int(b), 64)
				if err != nil {
					return nil, err
				}
				return states.NumValue(n), nil
			},
		),
	})
}
