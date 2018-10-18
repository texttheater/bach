package builtin

import (
	"math"

	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func Add(argFunctions []functions.Function) functions.Function {
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

func Subtract(argFunctions []functions.Function) functions.Function {
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

func Multiply(argFunctions []functions.Function) functions.Function {
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

func Divide(argFunctions []functions.Function) functions.Function {
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

func Modulo(argFunctions []functions.Function) functions.Function {
	return &functions.EvaluatorFunction{
		argFunctions,
		&types.NumberType{},
		func(inputValue values.Value, argumentValues []values.Value) values.Value {
			inputNumber := inputValue.(*values.NumberValue)
			argumentNumber := argumentValues[0].(*values.NumberValue)
			return &values.NumberValue{math.Mod(inputNumber.Value, argumentNumber.Value)}
		},
	}
}

func LessThan(argFunctions []functions.Function) functions.Function {
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

func GreaterThan(argFunctions []functions.Function) functions.Function {
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

func Equal(argFunctions []functions.Function) functions.Function {
	return &functions.EvaluatorFunction{
		argFunctions,
		&types.BooleanType{},
		func(inputValue values.Value, argumentValues []values.Value) values.Value {
			inputNumber := inputValue.(*values.NumberValue)
			argumentNumber := argumentValues[0].(*values.NumberValue)
			return &values.BooleanValue{inputNumber.Value == argumentNumber.Value}
		},
	}
}

func LessEqual(argFunctions []functions.Function) functions.Function {
	return &functions.EvaluatorFunction{
		argFunctions,
		&types.BooleanType{},
		func(inputValue values.Value, argumentValues []values.Value) values.Value {
			inputNumber := inputValue.(*values.NumberValue)
			argumentNumber := argumentValues[0].(*values.NumberValue)
			return &values.BooleanValue{inputNumber.Value <= argumentNumber.Value}
		},
	}
}

func GreaterEqual(argFunctions []functions.Function) functions.Function {
	return &functions.EvaluatorFunction{
		argFunctions,
		&types.BooleanType{},
		func(inputValue values.Value, argumentValues []values.Value) values.Value {
			inputNumber := inputValue.(*values.NumberValue)
			argumentNumber := argumentValues[0].(*values.NumberValue)
			return &values.BooleanValue{inputNumber.Value >= argumentNumber.Value}
		},
	}
}
