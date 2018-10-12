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

func And(argFunctions []shapes.Function) shapes.Function {
	return &functions.EvaluatorFunction{
		argFunctions,
		&types.BooleanType{},
		func(inputValue values.Value, argumentValues []values.Value) values.Value {
			inputBoolean := inputValue.(*values.BooleanValue)
			argumentBoolean := argumentValues[0].(*values.BooleanValue)
			return &values.BooleanValue{inputBoolean.Value && argumentBoolean.Value}
		},
	}
}

func Or(argFunctions []shapes.Function) shapes.Function {
	return &functions.EvaluatorFunction{
		argFunctions,
		&types.BooleanType{},
		func(inputValue values.Value, argumentValues []values.Value) values.Value {
			inputBoolean := inputValue.(*values.BooleanValue)
			argumentBoolean := argumentValues[0].(*values.BooleanValue)
			return &values.BooleanValue{inputBoolean.Value || argumentBoolean.Value}
		},
	}
}

func Not(argFunctions []shapes.Function) shapes.Function {
	return &functions.EvaluatorFunction{
		argFunctions,
		&types.BooleanType{},
		func(inputValue values.Value, argumentValues []values.Value) values.Value {
			inputBoolean := inputValue.(*values.BooleanValue)
			return &values.BooleanValue{!inputBoolean.Value}
		},
	}
}

func BooleanEqual(argFunctions []shapes.Function) shapes.Function {
	return &functions.EvaluatorFunction{
		argFunctions,
		&types.BooleanType{},
		func(inputValue values.Value, argumentValues []values.Value) values.Value {
			inputBoolean := inputValue.(*values.BooleanValue)
			argumentBoolean := argumentValues[0].(*values.BooleanValue)
			return &values.BooleanValue{inputBoolean.Value == argumentBoolean.Value}
		},
	}
}
