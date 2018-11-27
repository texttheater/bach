package functions

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type Function struct {
	InputType  types.Type
	Name       string
	Params     []*Parameter
	OutputType types.Type
	Action     Action
}

func (f Function) Signature() *Parameter {
	return &Parameter{
		InputType:  f.InputType,
		Name:       f.Name,
		Params:     f.Params,
		OutputType: f.OutputType,
	}
}

type Kernel func(inputValue values.Value, argValues []values.Value) values.Value

func SimpleFunction(inputType types.Type, name string, argTypes []types.Type,
	outputType types.Type, kernel Kernel) Function {
	params := make([]*Parameter, 0, len(argTypes))
	for _, argType := range argTypes {
		params = append(params, &Parameter{
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
		Action: func(inputState State, args []Action) State {
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
		},
	}
}
