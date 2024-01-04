package builtin

import (
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initNull() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.FuncerDefinition{
		expressions.SimpleFuncer(
			types.Any{},
			"null",
			nil,
			types.Null{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				return states.NullValue{}, nil
			},
		),
	})
}
