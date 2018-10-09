package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func Add(argFunctions []shapes.Function) shapes.Function {
	return &functions.EvaluatorFunction{
		argFunctions,
		&types.NumberType{},
		func(inputValue values.Value, argumentValues []values.Value) values.Value {
			inputNumber := inputValue.(*values.NumberValue)
			argumentNumber := argumentValues[0].(*values.NumberValue)
			return &values.NumberValue{inputNumber.Value + argumentNumber.Value}
		},
	}
}

func Subtract(argFunctions []shapes.Function) shapes.Function {
	return &functions.EvaluatorFunction{
		argFunctions,
		&types.NumberType{},
		func(inputValue values.Value, argumentValues []values.Value) values.Value {
			inputNumber := inputValue.(*values.NumberValue)
			argumentNumber := argumentValues[0].(*values.NumberValue)
			return &values.NumberValue{inputNumber.Value - argumentNumber.Value}
		},
	}
}

func Multiply(argFunctions []shapes.Function) shapes.Function {
	return &functions.EvaluatorFunction{
		argFunctions,
		&types.NumberType{},
		func(inputValue values.Value, argumentValues []values.Value) values.Value {
			inputNumber := inputValue.(*values.NumberValue)
			argumentNumber := argumentValues[0].(*values.NumberValue)
			return &values.NumberValue{inputNumber.Value * argumentNumber.Value}
		},
	}
}

func Divide(argFunctions []shapes.Function) shapes.Function {
	return &functions.EvaluatorFunction{
		argFunctions,
		&types.NumberType{},
		func(inputValue values.Value, argumentValues []values.Value) values.Value {
			inputNumber := inputValue.(*values.NumberValue)
			argumentNumber := argumentValues[0].(*values.NumberValue)
			return &values.NumberValue{inputNumber.Value / argumentNumber.Value}
		},
	}
}

func LessThan(argFunctions []shapes.Function) shapes.Function {
	return &functions.EvaluatorFunction{
		argFunctions,
		&types.BooleanType{},
		func(inputValue values.Value, argumentValues []values.Value) values.Value {
			inputNumber := inputValue.(*values.NumberValue)
			argumentNumber := argumentValues[0].(*values.NumberValue)
			return &values.BooleanValue{inputNumber.Value < argumentNumber.Value}
		},
	}
}

func GreaterThan(argFunctions []shapes.Function) shapes.Function {
	return &functions.EvaluatorFunction{
		argFunctions,
		&types.BooleanType{},
		func(inputValue values.Value, argumentValues []values.Value) values.Value {
			inputNumber := inputValue.(*values.NumberValue)
			argumentNumber := argumentValues[0].(*values.NumberValue)
			return &values.BooleanValue{inputNumber.Value > argumentNumber.Value}
		},
	}
}
