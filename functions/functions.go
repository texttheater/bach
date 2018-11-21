package functions

import (
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type Function struct {
	InputType  types.Type
	Name       string
	Params     []*parameters.Param
	OutputType types.Type
	Action     states.Action
}

func (f Function) Signature() *parameters.Param {
	return &parameters.Param{
		InputType:  f.InputType,
		Name:       f.Name,
		Params:     f.Params,
		OutputType: f.OutputType,
	}
}

type Kernel func(inputValue values.Value, argValues []values.Value) values.Value

func SimpleFunction(inputType types.Type, name string, argTypes []types.Type,
	outputType types.Type, kernel Kernel) Function {
	params := make([]*parameters.Param, 0, len(argTypes))
	for _, argType := range argTypes {
		params = append(params, &parameters.Param{
			InputType:  &types.AnyType{},
			Name:       "", // TODO ?
			Params:     nil,
			OutputType: argType,
		})
	}
	return Function{
		InputType:  inputType,
		Name:       name,
		Params:     params,
		OutputType: outputType,
		Action: states.Action{
			Name: name,
			Execute: func(inputState states.State, args []states.Action) states.State {
				argValues := make([]values.Value, 0, len(argTypes))
				argInputState := states.State{
					Value:       &values.NullValue{},
					ActionStack: inputState.ActionStack,
				}
				for _, arg := range args {
					argValues = append(argValues, arg.Execute(argInputState, nil).Value)
				}
				return states.State{
					Value:       kernel(inputState.Value, argValues),
					ActionStack: inputState.ActionStack,
				}
			},
		},
	}
}
