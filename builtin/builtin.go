package builtin

import (
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

var InitialShape = initialShape()

func initialShape() shapes.Shape {
	shape := shapes.Shape{&types.AnyType{}, nil}
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.NumberType{},
		"+",
		[]types.Type{&types.NumberType{}},
		func(argFunctions []shapes.Function) shapes.Function {
			return Add{argFunctions[0]}
		},
	})
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.NumberType{},
		"-",
		[]types.Type{&types.NumberType{}},
		func(argFunctions []shapes.Function) shapes.Function {
			return Subtract{argFunctions[0]}
		},
	})
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.NumberType{},
		"*",
		[]types.Type{&types.NumberType{}},
		func(argFunctions []shapes.Function) shapes.Function {
			return Multiply{argFunctions[0]}
		},
	})
	shape.Stack = shape.Stack.Push(shapes.NFF{
		&types.NumberType{},
		"/",
		[]types.Type{&types.NumberType{}},
		func(argFunctions []shapes.Function) shapes.Function {
			return Divide{argFunctions[0]}
		},
	})
	return shape
}

var InitialState = states.State{&values.NullValue{}, nil}
