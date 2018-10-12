package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func True(argFunctions []shapes.Function) shapes.Function {
	return &functions.EvaluatorFunction{
		argFunctions,
		&types.BooleanType{},
		func(inputValue values.Value, argumentValues []values.Value) values.Value {
			return &values.BooleanValue{true}
		},
	}
}

func False(argFunctions []shapes.Function) shapes.Function {
	return &functions.EvaluatorFunction{
		argFunctions,
		&types.BooleanType{},
		func(inputValue values.Value, argumentValues []values.Value) values.Value {
			return &values.BooleanValue{false}
		},
	}
}
