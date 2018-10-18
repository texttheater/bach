package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

var InitialShape = initialShape()

func initialShape() functions.Shape {
	shape := functions.Shape{&types.AnyType{}, nil}
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"+",
		[]types.Type{&types.NumberType{}},
		&types.NumberType{},
		Add,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"-",
		[]types.Type{&types.NumberType{}},
		&types.NumberType{},
		Subtract,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"*",
		[]types.Type{&types.NumberType{}},
		&types.NumberType{},
		Multiply,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"/",
		[]types.Type{&types.NumberType{}},
		&types.NumberType{},
		Divide,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"%",
		[]types.Type{&types.NumberType{}},
		&types.NumberType{},
		Modulo,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"<",
		[]types.Type{&types.NumberType{}},
		&types.BooleanType{},
		LessThan,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		">",
		[]types.Type{&types.NumberType{}},
		&types.BooleanType{},
		GreaterThan,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"==",
		[]types.Type{&types.NumberType{}},
		&types.BooleanType{},
		Equal,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"<=",
		[]types.Type{&types.NumberType{}},
		&types.BooleanType{},
		LessEqual,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		">=",
		[]types.Type{&types.NumberType{}},
		&types.BooleanType{},
		GreaterEqual,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.AnyType{},
		"true",
		[]types.Type{},
		&types.BooleanType{},
		True,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.AnyType{},
		"false",
		[]types.Type{},
		&types.BooleanType{},
		False,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.BooleanType{},
		"and",
		[]types.Type{&types.BooleanType{}},
		&types.BooleanType{},
		And,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.BooleanType{},
		"or",
		[]types.Type{&types.BooleanType{}},
		&types.BooleanType{},
		Or,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.BooleanType{},
		"not",
		[]types.Type{},
		&types.BooleanType{},
		Not,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.BooleanType{},
		"==",
		[]types.Type{&types.BooleanType{}},
		&types.BooleanType{},
		BooleanEqual,
	})
	return shape
}
