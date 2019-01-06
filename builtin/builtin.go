package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

var InitialShape = initialShape()

func initialShape() functions.Shape {
	shape := functions.Shape{&types.NullType{}, nil}
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		&types.NumType{},
		"+",
		[]types.Type{&types.NumType{}},
		&types.NumType{},
		Add,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		&types.NumType{},
		"-",
		[]types.Type{&types.NumType{}},
		&types.NumType{},
		Subtract,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		&types.NumType{},
		"*",
		[]types.Type{&types.NumType{}},
		&types.NumType{},
		Multiply,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		&types.NumType{},
		"/",
		[]types.Type{&types.NumType{}},
		&types.NumType{},
		Divide,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		&types.NumType{},
		"%",
		[]types.Type{&types.NumType{}},
		&types.NumType{},
		Modulo,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		&types.NumType{},
		"<",
		[]types.Type{&types.NumType{}},
		&types.BoolType{},
		LessThan,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		&types.NumType{},
		">",
		[]types.Type{&types.NumType{}},
		&types.BoolType{},
		GreaterThan,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		&types.NumType{},
		"==",
		[]types.Type{&types.NumType{}},
		&types.BoolType{},
		Equal,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		&types.NumType{},
		"<=",
		[]types.Type{&types.NumType{}},
		&types.BoolType{},
		LessEqual,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		&types.NumType{},
		">=",
		[]types.Type{&types.NumType{}},
		&types.BoolType{},
		GreaterEqual,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		&types.AnyType{},
		"true",
		nil,
		&types.BoolType{},
		True,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		&types.AnyType{},
		"false",
		nil,
		&types.BoolType{},
		False,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		&types.BoolType{},
		"and",
		[]types.Type{&types.BoolType{}},
		&types.BoolType{},
		And,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		&types.BoolType{},
		"or",
		[]types.Type{&types.BoolType{}},
		&types.BoolType{},
		Or,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		&types.BoolType{},
		"not",
		nil,
		&types.BoolType{},
		Not,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		&types.BoolType{},
		"==",
		[]types.Type{&types.BoolType{}},
		&types.BoolType{},
		BoolEqual,
	))
	shape.FuncerStack = shape.FuncerStack.Push(functions.SimpleFuncer(
		&types.AnyType{},
		"null",
		nil,
		&types.NullType{},
		Null,
	))
	return shape
}
