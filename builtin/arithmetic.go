package builtin

import (
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

// TODO need more abstraction below

///////////////////////////////////////////////////////////////////////////////

type Add struct {
	Arg shapes.Function
}

func (f Add) OutputShape(inputShape shapes.Shape) shapes.Shape {
	return inputShape
}

func (f Add) OutputState(inputState states.State) states.State {
	numberInput, _ := inputState.Value.(*values.NumberValue)
	argValue := f.Arg.OutputState(InitialState).Value
	numberArgValue, _ := argValue.(*values.NumberValue)
	outputValue := &values.NumberValue{numberInput.Value + numberArgValue.Value}
	return states.State{outputValue, inputState.Stack}
}

///////////////////////////////////////////////////////////////////////////////

type Subtract struct {
	Arg shapes.Function
}

func (f Subtract) OutputShape(inputShape shapes.Shape) shapes.Shape {
	return inputShape
}

func (f Subtract) OutputState(inputState states.State) states.State {
	numberInput, _ := inputState.Value.(*values.NumberValue)
	argValue := f.Arg.OutputState(InitialState).Value
	numberArgValue, _ := argValue.(*values.NumberValue)
	outputValue := &values.NumberValue{numberInput.Value - numberArgValue.Value}
	return states.State{outputValue, inputState.Stack}
}

///////////////////////////////////////////////////////////////////////////////

type Multiply struct {
	Arg shapes.Function
}

func (f Multiply) OutputShape(inputShape shapes.Shape) shapes.Shape {
	return inputShape
}

func (f Multiply) OutputState(inputState states.State) states.State {
	numberInput, _ := inputState.Value.(*values.NumberValue)
	argValue := f.Arg.OutputState(InitialState).Value
	numberArgValue, _ := argValue.(*values.NumberValue)
	outputValue := &values.NumberValue{numberInput.Value * numberArgValue.Value}
	return states.State{outputValue, inputState.Stack}
}

///////////////////////////////////////////////////////////////////////////////

type Divide struct {
	Arg shapes.Function
}

func (f Divide) OutputShape(inputShape shapes.Shape) shapes.Shape {
	return inputShape
}

func (f Divide) OutputState(inputState states.State) states.State {
	numberInput, _ := inputState.Value.(*values.NumberValue)
	argValue := f.Arg.OutputState(InitialState).Value
	numberArgValue, _ := argValue.(*values.NumberValue)
	outputValue := &values.NumberValue{numberInput.Value / numberArgValue.Value}
	return states.State{outputValue, inputState.Stack}
}

///////////////////////////////////////////////////////////////////////////////

type LessThan struct {
	Arg shapes.Function
}

func (f LessThan) OutputShape(inputShape shapes.Shape) shapes.Shape {
	return shapes.Shape{&types.BooleanType{}, inputShape.Stack}
}

func (f LessThan) OutputState(inputState states.State) states.State {
	numberInput, _ := inputState.Value.(*values.NumberValue)
	argValue := f.Arg.OutputState(InitialState).Value
	numberArgValue, _ := argValue.(*values.NumberValue)
	outputValue := &values.BooleanValue{numberInput.Value < numberArgValue.Value}
	return states.State{outputValue, inputState.Stack}
}

///////////////////////////////////////////////////////////////////////////////

type GreaterThan struct {
	Arg shapes.Function
}

func (f GreaterThan) OutputShape(inputShape shapes.Shape) shapes.Shape {
	return shapes.Shape{&types.BooleanType{}, inputShape.Stack}
}

func (f GreaterThan) OutputState(inputState states.State) states.State {
	numberInput, _ := inputState.Value.(*values.NumberValue)
	argValue := f.Arg.OutputState(InitialState).Value
	numberArgValue, _ := argValue.(*values.NumberValue)
	outputValue := &values.BooleanValue{numberInput.Value > numberArgValue.Value}
	return states.State{outputValue, inputState.Stack}
}

///////////////////////////////////////////////////////////////////////////////
