package builtin

import (
	"github.com/texttheater/bach/values"
)

func True(inputValue values.Value, argumentValues []values.Value) values.Value {
	return &values.BooleanValue{true}
}

func False(inputValue values.Value, argumentValues []values.Value) values.Value {
	return &values.BooleanValue{false}
}

func And(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputBoolean := inputValue.(*values.BooleanValue)
	argumentBoolean := argumentValues[0].(*values.BooleanValue)
	return &values.BooleanValue{inputBoolean.Value && argumentBoolean.Value}
}

func Or(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputBoolean := inputValue.(*values.BooleanValue)
	argumentBoolean := argumentValues[0].(*values.BooleanValue)
	return &values.BooleanValue{inputBoolean.Value || argumentBoolean.Value}
}

func Not(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputBoolean := inputValue.(*values.BooleanValue)
	return &values.BooleanValue{!inputBoolean.Value}
}

func BooleanEqual(inputValue values.Value, argumentValues []values.Value) values.Value {
	inputBoolean := inputValue.(*values.BooleanValue)
	argumentBoolean := argumentValues[0].(*values.BooleanValue)
	return &values.BooleanValue{inputBoolean.Value == argumentBoolean.Value}
}
