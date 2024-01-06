package builtin

import (
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var NullFuncers = []shapes.Funcer{
	shapes.SimpleFuncer(
		types.Any{},
		"null",
		nil,
		types.Null{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			return states.NullValue{}, nil
		},
	),
}
