package functions

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type Function struct {
	InputType  types.Type
	Name       string
	Params     []*Param
	OutputType types.Type
	Action     Action
}

type Kernel func(inputValue values.Value, argValues []values.Value) values.Value

func SimpleFunction(inputType types.Type, name string, argTypes []types.Type,
	outputType types.Type, kernel Kernel) Function {
	params := make([]*Param, 0, len(argTypes))
	for _, argType := range argTypes {
		params = append(params, &Param{
			InputType:  &types.NullType{},
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
		Action: func(inputValue values.Value, args []Action) values.Value {
			argValues := make([]values.Value, 0, len(argTypes))
			for _, arg := range args {
				argValues = append(argValues, arg(&values.NullValue{}, nil))
			}
			return kernel(inputValue, argValues)
		},
	}
}
