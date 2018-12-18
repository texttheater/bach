package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

var InitialShape = initialShape()

func initialShape() functions.Shape {
	shape := functions.Shape{&types.NullType{}, nil}
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumType{},
		"+",
		[]types.Type{&types.NumType{}},
		&types.NumType{},
		Add,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumType{},
		"-",
		[]types.Type{&types.NumType{}},
		&types.NumType{},
		Subtract,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumType{},
		"*",
		[]types.Type{&types.NumType{}},
		&types.NumType{},
		Multiply,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumType{},
		"/",
		[]types.Type{&types.NumType{}},
		&types.NumType{},
		Divide,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumType{},
		"%",
		[]types.Type{&types.NumType{}},
		&types.NumType{},
		Modulo,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumType{},
		"<",
		[]types.Type{&types.NumType{}},
		&types.BoolType{},
		LessThan,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumType{},
		">",
		[]types.Type{&types.NumType{}},
		&types.BoolType{},
		GreaterThan,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumType{},
		"==",
		[]types.Type{&types.NumType{}},
		&types.BoolType{},
		Equal,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumType{},
		"<=",
		[]types.Type{&types.NumType{}},
		&types.BoolType{},
		LessEqual,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.NumType{},
		">=",
		[]types.Type{&types.NumType{}},
		&types.BoolType{},
		GreaterEqual,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.AnyType{},
		"true",
		nil,
		&types.BoolType{},
		True,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.AnyType{},
		"false",
		nil,
		&types.BoolType{},
		False,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.BoolType{},
		"and",
		[]types.Type{&types.BoolType{}},
		&types.BoolType{},
		And,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.BoolType{},
		"or",
		[]types.Type{&types.BoolType{}},
		&types.BoolType{},
		Or,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.BoolType{},
		"not",
		nil,
		&types.BoolType{},
		Not,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.BoolType{},
		"==",
		[]types.Type{&types.BoolType{}},
		&types.BoolType{},
		BoolEqual,
	))
	shape.FunctionStack = shape.FunctionStack.Push(functions.SimpleFunction(
		&types.AnyType{},
		"null",
		nil,
		&types.NullType{},
		Null,
	))
	return shape
}
