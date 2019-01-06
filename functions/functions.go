package functions

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type Funcer func(gotInputType types.Type, gotName string, gotNumArgs int) (params []*Parameter, outputType types.Type, action Action, ok bool)

type Kernel func(inputValue values.Value, argValues []values.Value) values.Value

func SimpleFuncer(wantInputType types.Type, wantName string, argTypes []types.Type, outputType types.Type, kernel Kernel) Funcer {
	// make parameters from argument types
	params := make([]*Parameter, 0, len(argTypes))
	for _, argType := range argTypes {
		params = append(params, &Parameter{
			InputType:  &types.AnyType{},
			Name:       "", // TODO ?
			Params:     nil,
			OutputType: argType,
		})
	}
	// make funcer
	return func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*Parameter, types.Type, Action, bool) {
		if !wantInputType.Subsumes(gotInputType) {
			return nil, nil, nil, false
		}
		if gotName != wantName {
			return nil, nil, nil, false
		}
		if gotNumArgs != len(argTypes) {
			return nil, nil, nil, false
		}
		action := func(inputState State, args []Action) State {
			argValues := make([]values.Value, 0, len(argTypes))
			argInputState := State{
				Value: &values.NullValue{},
				Stack: inputState.Stack,
			}
			for _, arg := range args {
				argValues = append(argValues, arg(argInputState, nil).Value)
			}
			return State{
				Value: kernel(inputState.Value, argValues),
				Stack: inputState.Stack,
			}
		}
		return params, outputType, action, true
	}
}
