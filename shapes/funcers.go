package shapes

import (
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type Funcer func(gotInputType types.Type, gotName string, gotNumArgs int) (params []*parameters.Parameter, outputType types.Type, action states.Action, ok bool)

type Kernel func(inputValue values.Value, argValues []values.Value) values.Value

func SimpleFuncer(wantInputType types.Type, wantName string, argTypes []types.Type, outputType types.Type, kernel Kernel) Funcer {
	// make parameters from argument types
	params := make([]*parameters.Parameter, len(argTypes))
	for i, argType := range argTypes {
		params[i] = &parameters.Parameter{
			InputType:  types.AnyType,
			Name:       "", // TODO ?
			Params:     nil,
			OutputType: argType,
		}
	}
	// make funcer
	return func(gotInputType types.Type, gotName string, gotNumArgs int) ([]*parameters.Parameter, types.Type, states.Action, bool) {
		if !wantInputType.Subsumes(gotInputType) {
			return nil, nil, nil, false
		}
		if gotName != wantName {
			return nil, nil, nil, false
		}
		if gotNumArgs != len(argTypes) {
			return nil, nil, nil, false
		}
		action := func(inputState states.State, args []states.Action) states.State {
			argValues := make([]values.Value, len(argTypes))
			argInputState := states.State{
				Value: &values.NullValue{},
				Stack: inputState.Stack,
			}
			for i, arg := range args {
				argValues[i] = arg(argInputState, nil).Value
			}
			return states.State{
				Value: kernel(inputState.Value, argValues),
				Stack: inputState.Stack,
			}
		}
		return params, outputType, action, true
	}
}
