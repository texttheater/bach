package shapes

import (
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/values"
)

// TODO need more abstraction below

///////////////////////////////////////////////////////////////////////////////

type Add struct {
	Arg Function
}

func (f Add) OutputShape(inputShape Shape) Shape {
	return inputShape
}

func (f Add) OutputState(inputState states.State) states.State {
	numberInput, _ := inputState.Value.(*values.NumberValue)
	argValue := f.Arg.OutputState(states.InitialState).Value
	numberArgValue, _ := argValue.(*values.NumberValue)
	outputValue := &values.NumberValue{numberInput.Value + numberArgValue.Value}
	return states.State{outputValue, inputState.Stack}
}

///////////////////////////////////////////////////////////////////////////////

type Subtract struct {
	Arg Function
}

func (f Subtract) OutputShape(inputShape Shape) Shape {
	return inputShape
}

func (f Subtract) OutputState(inputState states.State) states.State {
	numberInput, _ := inputState.Value.(*values.NumberValue)
	argValue := f.Arg.OutputState(states.InitialState).Value
	numberArgValue, _ := argValue.(*values.NumberValue)
	outputValue := &values.NumberValue{numberInput.Value - numberArgValue.Value}
	return states.State{outputValue, inputState.Stack}
}

///////////////////////////////////////////////////////////////////////////////

type Multiply struct {
	Arg Function
}

func (f Multiply) OutputShape(inputShape Shape) Shape {
	return inputShape
}

func (f Multiply) OutputState(inputState states.State) states.State {
	numberInput, _ := inputState.Value.(*values.NumberValue)
	argValue := f.Arg.OutputState(states.InitialState).Value
	numberArgValue, _ := argValue.(*values.NumberValue)
	outputValue := &values.NumberValue{numberInput.Value * numberArgValue.Value}
	return states.State{outputValue, inputState.Stack}
}

///////////////////////////////////////////////////////////////////////////////

type Divide struct {
	Arg Function
}

func (f Divide) OutputShape(inputShape Shape) Shape {
	return inputShape
}

func (f Divide) OutputState(inputState states.State) states.State {
	numberInput, _ := inputState.Value.(*values.NumberValue)
	argValue := f.Arg.OutputState(states.InitialState).Value
	numberArgValue, _ := argValue.(*values.NumberValue)
	outputValue := &values.NumberValue{numberInput.Value / numberArgValue.Value}
	return states.State{outputValue, inputState.Stack}
}

///////////////////////////////////////////////////////////////////////////////
