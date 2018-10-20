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
		functions.DumbFuncer(&types.NumberType{}, Add),
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"-",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		functions.DumbFuncer(&types.NumberType{}, Subtract),
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"*",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		functions.DumbFuncer(&types.NumberType{}, Multiply),
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"/",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		functions.DumbFuncer(&types.NumberType{}, Divide),
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"%",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		functions.DumbFuncer(&types.NumberType{}, Modulo),
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"<",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		functions.DumbFuncer(&types.BooleanType{}, LessThan),
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		">",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		functions.DumbFuncer(&types.BooleanType{}, GreaterThan),
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"==",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		functions.DumbFuncer(&types.BooleanType{}, Equal),
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"<=",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		functions.DumbFuncer(&types.BooleanType{}, LessEqual),
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		">=",
		[]*parameters.Parameter{parameters.DumbPar(&types.NumberType{})},
		functions.DumbFuncer(&types.BooleanType{}, GreaterEqual),
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.AnyType{},
		"true",
		[]*parameters.Parameter{},
		functions.DumbFuncer(&types.BooleanType{}, True),
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.AnyType{},
		"false",
		[]*parameters.Parameter{},
		functions.DumbFuncer(&types.BooleanType{}, False),
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.BooleanType{},
		"and",
		[]*parameters.Parameter{parameters.DumbPar(&types.BooleanType{})},
		functions.DumbFuncer(&types.BooleanType{}, And),
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.BooleanType{},
		"or",
		[]*parameters.Parameter{parameters.DumbPar(&types.BooleanType{})},
		functions.DumbFuncer(&types.BooleanType{}, Or),
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.BooleanType{},
		"not",
		[]*parameters.Parameter{},
		functions.DumbFuncer(&types.BooleanType{}, Not),
	})
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.BooleanType{},
		"==",
		[]*parameters.Parameter{parameters.DumbPar(&types.BooleanType{})},
		functions.DumbFuncer(&types.BooleanType{}, BooleanEqual),
	})
	// TODO remove; this is just for testing
	shape.Stack = shape.Stack.Push(functions.NFF{
		&types.NumberType{},
		"apply",
		[]*parameters.Parameter{
			&parameters.Parameter{
				&types.NumberType{},
				[]*parameters.Parameter{},
				&types.NumberType{},
			},
		},
		func(argFuncers ...functions.Funcer) functions.Function {
			argFunction := argFuncers[0]()
			return &functions.ApplyFunction{argFunction}
		},
	})
	//	shape.Stack = shape.Stack.Push(functions.NFF{
	//		&types.AnyType{},
	//		"one_o_one",
	//		[]*parameters.Parameter{
	//			&parameters.Parameter{
	//				&types.NumberType{},
	//				[]*parameters.Parameter{
	//					&parameters.Parameter{
	//						&types.AnyType{},
	//						[]*parameters.Parameter{},
	//						&types.NumberType{},
	//					},
	//				},
	//				&types.NumberType{},
	//			}
	//		},
	//		func(argFunctions []functions.Function) functions.Function {
	//
	//		}
	//	})
	return shape
}
