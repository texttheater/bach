/*
Package shapes implements shapes.

A shape consists of a type and a stack of available NFFs.

Interpreting a Bach program involves assigning each expression an input shape,
a function and an output shape. The first expression in the program gets the
initial shape, consisting of the Any type and a stack consisting only of
builtin NFFs. The input shape of an expression and the expression together
determine its function. The function and the input shape together determine
its output shape. In a concatenation expression L R, the output shape of L is
the input shape of R.
*/
// TODO tear this package apart
package shapes

import (
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type Shape struct {
	Type types.Type
	Stack *Stack
}

type Stack struct {
	Head NFF
	Tail *Stack
}

func (stack *Stack) Push(n NFF) *Stack {
	return &Stack{n, stack}
}

func (stack *Stack) Pop() *Stack {
	return stack.Tail
}

type NFF struct {
	InputType types.Type // TODO type parameters
	Name string // TODO namespaces
	ArgTypes []types.Type
	Funcer func([]Function) Function // TODO first-class functions
}

type Function interface {
	OutputShape(inputShape Shape) Shape
	OutputState(inputState states.State) states.State
}

var InitialShape = initialShape()

func initialShape() Shape {
	shape := Shape{&types.AnyType{}, nil}
	shape.Stack = shape.Stack.Push(NFF{
		&types.NumberType{},
		"+",
		[]types.Type{&types.NumberType{}},
		func(argFunctions []Function) Function {
			return Add{argFunctions[0]}
		},
	})
	shape.Stack = shape.Stack.Push(NFF{
		&types.NumberType{},
		"-",
		[]types.Type{&types.NumberType{}},
		func(argFunctions []Function) Function {
			return Subtract{argFunctions[0]}
		},
	})
	shape.Stack = shape.Stack.Push(NFF{
		&types.NumberType{},
		"*",
		[]types.Type{&types.NumberType{}},
		func(argFunctions []Function) Function {
			return Multiply{argFunctions[0]}
		},
	})
	shape.Stack = shape.Stack.Push(NFF{
		&types.NumberType{},
		"/",
		[]types.Type{&types.NumberType{}},
		func(argFunctions []Function) Function {
			return Divide{argFunctions[0]}
		},
	})
	return shape
}
