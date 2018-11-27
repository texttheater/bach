package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

var InitialShape = initialShape()

func initialShape() functions.Shape {
	shape := functions.Shape{&types.AnyType{}, nil}
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"+",
		[]types.Type{&types.NumberType{}},
		&types.NumberType{},
		Add,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"-",
		[]types.Type{&types.NumberType{}},
		&types.NumberType{},
		Subtract,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"*",
		[]types.Type{&types.NumberType{}},
		&types.NumberType{},
		Multiply,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"/",
		[]types.Type{&types.NumberType{}},
		&types.NumberType{},
		Divide,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"%",
		[]types.Type{&types.NumberType{}},
		&types.NumberType{},
		Modulo,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"<",
		[]types.Type{&types.NumberType{}},
		&types.BooleanType{},
		LessThan,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		">",
		[]types.Type{&types.NumberType{}},
		&types.BooleanType{},
		GreaterThan,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"==",
		[]types.Type{&types.NumberType{}},
		&types.BooleanType{},
		Equal,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"<=",
		[]types.Type{&types.NumberType{}},
		&types.BooleanType{},
		LessEqual,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		">=",
		[]types.Type{&types.NumberType{}},
		&types.BooleanType{},
		GreaterEqual,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.AnyType{},
		"true",
		nil,
		&types.BooleanType{},
		True,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.AnyType{},
		"false",
		nil,
		&types.BooleanType{},
		False,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.BooleanType{},
		"and",
		[]types.Type{&types.BooleanType{}},
		&types.BooleanType{},
		And,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.BooleanType{},
		"or",
		[]types.Type{&types.BooleanType{}},
		&types.BooleanType{},
		Or,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.BooleanType{},
		"not",
		nil,
		&types.BooleanType{},
		Not,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.BooleanType{},
		"==",
		[]types.Type{&types.BooleanType{}},
		&types.BooleanType{},
		BooleanEqual,
	))
	return shape
}
