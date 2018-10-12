package builtin

import (
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/types"
)

var InitialShape = initialShape()

func initialShape() shapes.Shape {
	shape := shapes.Shape{&types.AnyType{}, nil}
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.NumberType{},
		"+",
		[]types.Type{&types.NumberType{}},
		Add,
	})
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.NumberType{},
		"-",
		[]types.Type{&types.NumberType{}},
		Subtract,
	})
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.NumberType{},
		"*",
		[]types.Type{&types.NumberType{}},
		Multiply,
	})
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.NumberType{},
		"/",
		[]types.Type{&types.NumberType{}},
		Divide,
	})
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.NumberType{},
		"%",
		[]types.Type{&types.NumberType{}},
		Modulo,
	})
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.NumberType{},
		"<",
		[]types.Type{&types.NumberType{}},
		LessThan,
	})
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.NumberType{},
		">",
		[]types.Type{&types.NumberType{}},
		GreaterThan,
	})
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.NumberType{},
		"==",
		[]types.Type{&types.NumberType{}},
		Equal,
	})
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.NumberType{},
		"<=",
		[]types.Type{&types.NumberType{}},
		LessEqual,
	})
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.NumberType{},
		">=",
		[]types.Type{&types.NumberType{}},
		GreaterEqual,
	})
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.AnyType{},
		"true",
		[]types.Type{},
		True,
	})
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.AnyType{},
		"false",
		[]types.Type{},
		False,
	})
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.BooleanType{},
		"and",
		[]types.Type{&types.BooleanType{}},
		And,
	})
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.BooleanType{},
		"or",
		[]types.Type{&types.BooleanType{}},
		Or,
	})
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.BooleanType{},
		"not",
		[]types.Type{},
		Not,
	})
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.BooleanType{},
		"==",
		[]types.Type{&types.BooleanType{}},
		BooleanEqual,
	})
	return shape
}
