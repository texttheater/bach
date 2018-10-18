// Package functions implements functions from states to states.
//
// A state consists of a value and a stack of available variables (named
// values).
//
// Shapes are to states as types are to values. A shape consists of a type and
// a stack of available NFFs (named function families).
//
// Interpreting a Bach program involves assigning each expression an input
// shape, a function and an output shape. The first expression in the program
// gets the initial shape, consisting of the Any type and a stack consisting
// only of builtin NFFs. The input shape of an expression and the expression
// together determine its function. The function and the input shape together
// determine its output shape. In a composition expression L R, the output
// shape of L is the input shape of R.
package functions

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

///////////////////////////////////////////////////////////////////////////////

type Function interface {
	OutputShape(inputStack *NFFStack) Shape
	OutputState(inputState State) State
}

///////////////////////////////////////////////////////////////////////////////

type IdentityFunction struct {
	Type types.Type
}

func (f *IdentityFunction) OutputShape(inputStack *NFFStack) Shape {
	return Shape{f.Type, inputStack}
}

func (f *IdentityFunction) OutputState(inputState State) State {
	return inputState
}

///////////////////////////////////////////////////////////////////////////////

type CompositionFunction struct {
	Left  Function
	Right Function
}

func (f *CompositionFunction) OutputShape(inputStack *NFFStack) Shape {
	return f.Right.OutputShape(f.Left.OutputShape(inputStack).Stack)
}

func (f *CompositionFunction) OutputState(inputState State) State {
	return f.Right.OutputState(f.Left.OutputState(inputState))
}

///////////////////////////////////////////////////////////////////////////////

type LiteralFunction struct {
	Type  types.Type
	Value values.Value
}

func (f *LiteralFunction) OutputShape(inputStack *NFFStack) Shape {
	return Shape{f.Type, inputStack}
}

func (f *LiteralFunction) OutputState(inputState State) State {
	return State{f.Value, inputState.Stack}
}

///////////////////////////////////////////////////////////////////////////////

type EvaluatorFunction struct {
	ArgumentFunctions []Function
	OutputType        types.Type
	Kernel            Kernel
}

func (f *EvaluatorFunction) OutputShape(inputStack *NFFStack) Shape {
	return Shape{f.OutputType, inputStack}
}

func (f *EvaluatorFunction) OutputState(inputState State) State {
	argumentInput := State{&values.NullValue{}, inputState.Stack}
	argumentValues := make([]values.Value, len(f.ArgumentFunctions))
	for i, a := range f.ArgumentFunctions {
		argumentValues[i] = a.OutputState(argumentInput).Value
	}
	return State{f.Kernel(inputState, argumentValues), inputState.Stack}
}

///////////////////////////////////////////////////////////////////////////////

type AssignmentFunction struct {
	Type types.Type
	Name string
}

func (f *AssignmentFunction) OutputShape(inputStack *NFFStack) Shape {
	return Shape{f.Type, inputStack.Push(NFF{
		&types.AnyType{},
		f.Name,
		[]types.Type{},
		f.Type,
		func(inputState State, argumentValues []values.Value) values.Value {
			stack := inputState.Stack
			for stack != nil {
				if stack.Head.Name == f.Name {
					return stack.Head.Value
				}
				stack = stack.Tail
			}
			panic("variable not found")
		},
	})}
}

func (f *AssignmentFunction) OutputState(inputState State) State {
	return State{
		inputState.Value,
		inputState.Stack.Push(NamedValue{
			f.Name,
			inputState.Value,
		}),
	}
}

///////////////////////////////////////////////////////////////////////////////
