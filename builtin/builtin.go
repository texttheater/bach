package builtin

import (
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

///////////////////////////////////////////////////////////////////////////////

type EvaluatorFunction struct {
	ArgumentFunctions []shapes.Function
	OutputType        types.Type
	Kernel            func(inputValue values.Value, argumentValues []values.Value) values.Value
}

func (f *EvaluatorFunction) OutputShape(inputShape shapes.Shape) shapes.Shape {
	return shapes.Shape{f.OutputType, inputShape.Stack}
}

func (f *EvaluatorFunction) OutputState(inputState states.State) states.State {
	argumentValues := make([]values.Value, len(f.ArgumentFunctions))
	for i, a := range f.ArgumentFunctions {
		argumentValues[i] = a.OutputState(InitialState).Value
	}
	return states.State{f.Kernel(inputState.Value, argumentValues), inputState.Stack}
}

///////////////////////////////////////////////////////////////////////////////

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
	return shape
}

///////////////////////////////////////////////////////////////////////////////

var InitialState = states.State{&values.NullValue{}, nil}
