package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/values"
)

func True(inputState functions.State, argumentValues []values.Value) values.Value {
	return &values.BooleanValue{true}
}

func False(inputState functions.State, argumentValues []values.Value) values.Value {
	return &values.BooleanValue{false}
}

func And(inputState functions.State, argumentValues []values.Value) values.Value {
	inputBoolean := inputState.Value.(*values.BooleanValue)
	argumentBoolean := argumentValues[0].(*values.BooleanValue)
	return &values.BooleanValue{inputBoolean.Value && argumentBoolean.Value}
}

func Or(inputState functions.State, argumentValues []values.Value) values.Value {
	inputBoolean := inputState.Value.(*values.BooleanValue)
	argumentBoolean := argumentValues[0].(*values.BooleanValue)
	return &values.BooleanValue{inputBoolean.Value || argumentBoolean.Value}
}

func Not(inputState functions.State, argumentValues []values.Value) values.Value {
	inputBoolean := inputState.Value.(*values.BooleanValue)
	return &values.BooleanValue{!inputBoolean.Value}
}

func BooleanEqual(inputState functions.State, argumentValues []values.Value) values.Value {
	inputBoolean := inputState.Value.(*values.BooleanValue)
	argumentBoolean := argumentValues[0].(*values.BooleanValue)
	return &values.BooleanValue{inputBoolean.Value == argumentBoolean.Value}
}
