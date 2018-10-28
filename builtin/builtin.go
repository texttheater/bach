package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

var InitialShape = initialShape()

func initialShape() functions.Shape {
	shape := functions.Shape{&types.AnyType{}, nil}
	shape.Stack = shape.Stack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"+",
		[]*functions.Parameter{functions.SimpleParameter("summand", &types.NumberType{})},
		&types.NumberType{},
		Add,
	))
	shape.Stack = shape.Stack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"-",
		[]*functions.Parameter{functions.SimpleParameter("subtrahend", &types.NumberType{})},
		&types.NumberType{},
		Subtract,
	))
	shape.Stack = shape.Stack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"*",
		[]*functions.Parameter{functions.SimpleParameter("factor", &types.NumberType{})},
		&types.NumberType{},
		Multiply,
	))
	shape.Stack = shape.Stack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"/",
		[]*functions.Parameter{functions.SimpleParameter("divisor", &types.NumberType{})},
		&types.NumberType{},
		Divide,
	))
	shape.Stack = shape.Stack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"%",
		[]*functions.Parameter{functions.SimpleParameter("divisor", &types.NumberType{})},
		&types.NumberType{},
		Modulo,
	))
	shape.Stack = shape.Stack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"<",
		[]*functions.Parameter{functions.SimpleParameter("other", &types.NumberType{})},
		&types.BooleanType{},
		LessThan,
	))
	shape.Stack = shape.Stack.Push(functions.SimpleFunction(
		&types.NumberType{},
		">",
		[]*functions.Parameter{functions.SimpleParameter("other", &types.NumberType{})},
		&types.BooleanType{},
		GreaterThan,
	))
	shape.Stack = shape.Stack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"==",
		[]*functions.Parameter{functions.SimpleParameter("other", &types.NumberType{})},
		&types.BooleanType{},
		Equal,
	))
	shape.Stack = shape.Stack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"<=",
		[]*functions.Parameter{functions.SimpleParameter("other", &types.NumberType{})},
		&types.BooleanType{},
		LessEqual,
	))
	shape.Stack = shape.Stack.Push(functions.SimpleFunction(
		&types.NumberType{},
		">=",
		[]*functions.Parameter{functions.SimpleParameter("other", &types.NumberType{})},
		&types.BooleanType{},
		GreaterEqual,
	))
	shape.Stack = shape.Stack.Push(functions.SimpleFunction(
		&types.AnyType{},
		"true",
		nil,
		&types.BooleanType{},
		True,
	))
	shape.Stack = shape.Stack.Push(functions.SimpleFunction(
		&types.AnyType{},
		"false",
		nil,
		&types.BooleanType{},
		False,
	))
	shape.Stack = shape.Stack.Push(functions.SimpleFunction(
		&types.BooleanType{},
		"and",
		[]*functions.Parameter{functions.SimpleParameter("other", &types.BooleanType{})},
		&types.BooleanType{},
		And,
	))
	shape.Stack = shape.Stack.Push(functions.SimpleFunction(
		&types.BooleanType{},
		"or",
		[]*functions.Parameter{functions.SimpleParameter("other", &types.BooleanType{})},
		&types.BooleanType{},
		Or,
	))
	shape.Stack = shape.Stack.Push(functions.SimpleFunction(
		&types.BooleanType{},
		"not",
		nil,
		&types.BooleanType{},
		Not,
	))
	shape.Stack = shape.Stack.Push(functions.SimpleFunction(
		&types.BooleanType{},
		"==",
		[]*functions.Parameter{functions.SimpleParameter("other", &types.BooleanType{})},
		&types.BooleanType{},
		BooleanEqual,
	))
	return shape
}
