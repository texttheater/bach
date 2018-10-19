package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/types"
)

var InitialShape = initialShape()

func initialShape() functions.Shape {
	shape := functions.Shape{&types.AnyType{}, nil}
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"+",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		&types.NumberType{},
		Add,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"-",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		&types.NumberType{},
		Subtract,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"*",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		&types.NumberType{},
		Multiply,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"/",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		&types.NumberType{},
		Divide,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"%",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		&types.NumberType{},
		Modulo,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"<",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		&types.BooleanType{},
		LessThan,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		">",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		&types.BooleanType{},
		GreaterThan,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"==",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		&types.BooleanType{},
		Equal,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"<=",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		&types.BooleanType{},
		LessEqual,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		">=",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		&types.BooleanType{},
		GreaterEqual,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.AnyType{},
		"true",
		[]*parameters.Parameter{},
		&types.BooleanType{},
		True,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.AnyType{},
		"false",
		[]*parameters.Parameter{},
		&types.BooleanType{},
		False,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.BooleanType{},
		"and",
		[]*parameters.Parameter{parameters.DumbPar(&types.BooleanType{})},
		&types.BooleanType{},
		And,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.BooleanType{},
		"or",
		[]*parameters.Parameter{parameters.DumbPar(&types.BooleanType{})},
		&types.BooleanType{},
		Or,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.BooleanType{},
		"not",
		[]*parameters.Parameter{},
		&types.BooleanType{},
		Not,
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.BooleanType{},
		"==",
		[]*parameters.Parameter{parameters.DumbPar(&types.BooleanType{})},
		&types.BooleanType{},
		BooleanEqual,
	})
	return shape
}
