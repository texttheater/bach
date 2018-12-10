package builtin

import (
	"github.com/texttheater/bach/values"
)

func True(inputValue values.Value, argumentValues []values.Value) values.Value {
	return &values.BoolValue{true}
}

func False(inputValue values.Value, argumentValues []values.Value) values.Value {
	return &values.BoolValue{false}
}

func And(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputBool := inputValue.(*values.BoolValue)
	argumentBool := argumentValues[0].(*values.BoolValue)
	return &values.BoolValue{inputBool.Value && argumentBool.Value}
}

func Or(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputBool := inputValue.(*values.BoolValue)
	argumentBool := argumentValues[0].(*values.BoolValue)
	return &values.BoolValue{inputBool.Value || argumentBool.Value}
}

func Not(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputBool := inputValue.(*values.BoolValue)
	return &values.BoolValue{!inputBool.Value}
}

func BoolEqual(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputBool := inputValue.(*values.BoolValue)
	argumentBool := argumentValues[0].(*values.BoolValue)
	return &values.BoolValue{inputBool.Value == argumentBool.Value}
}
