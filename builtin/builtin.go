package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

var InitialContext = initialContext()

func initialContext() functions.Context {
	context := functions.Context{&types.AnyType{}, nil}
	context.FunctionStack = context.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"+",
		[]types.Type{&types.NumberType{}},
		&types.NumberType{},
		Add,
	))
	context.FunctionStack = context.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"-",
		[]types.Type{&types.NumberType{}},
		&types.NumberType{},
		Subtract,
	))
	context.FunctionStack = context.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"*",
		[]types.Type{&types.NumberType{}},
		&types.NumberType{},
		Multiply,
	))
	context.FunctionStack = context.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"/",
		[]types.Type{&types.NumberType{}},
		&types.NumberType{},
		Divide,
	))
	context.FunctionStack = context.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"%",
		[]types.Type{&types.NumberType{}},
		&types.NumberType{},
		Modulo,
	))
	context.FunctionStack = context.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"<",
		[]types.Type{&types.NumberType{}},
		&types.BooleanType{},
		LessThan,
	))
	context.FunctionStack = context.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		">",
		[]types.Type{&types.NumberType{}},
		&types.BooleanType{},
		GreaterThan,
	))
	context.FunctionStack = context.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"==",
		[]types.Type{&types.NumberType{}},
		&types.BooleanType{},
		Equal,
	))
	context.FunctionStack = context.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		"<=",
		[]types.Type{&types.NumberType{}},
		&types.BooleanType{},
		LessEqual,
	))
	context.FunctionStack = context.FunctionStack.Push(functions.SimpleFunction(
		&types.NumberType{},
		">=",
		[]types.Type{&types.NumberType{}},
		&types.BooleanType{},
		GreaterEqual,
	))
	context.FunctionStack = context.FunctionStack.Push(functions.SimpleFunction(
		&types.AnyType{},
		"true",
		nil,
		&types.BooleanType{},
		True,
	))
	context.FunctionStack = context.FunctionStack.Push(functions.SimpleFunction(
		&types.AnyType{},
		"false",
		nil,
		&types.BooleanType{},
		False,
	))
	context.FunctionStack = context.FunctionStack.Push(functions.SimpleFunction(
		&types.BooleanType{},
		"and",
		[]types.Type{&types.BooleanType{}},
		&types.BooleanType{},
		And,
	))
	context.FunctionStack = context.FunctionStack.Push(functions.SimpleFunction(
		&types.BooleanType{},
		"or",
		[]types.Type{&types.BooleanType{}},
		&types.BooleanType{},
		Or,
	))
	context.FunctionStack = context.FunctionStack.Push(functions.SimpleFunction(
		&types.BooleanType{},
		"not",
		nil,
		&types.BooleanType{},
		Not,
	))
	context.FunctionStack = context.FunctionStack.Push(functions.SimpleFunction(
		&types.BooleanType{},
		"==",
		[]types.Type{&types.BooleanType{}},
		&types.BooleanType{},
		BooleanEqual,
	))
	return context
}
