package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initNull() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.SimpleFuncer(
			types.AnyType{},
			"null",
			nil,
			types.NullType{},
			func(inputValue values.Value, argumentValues []values.Value) (values.Value, error) {
				return values.NullValue{}, nil
			},
		),
	})
}
