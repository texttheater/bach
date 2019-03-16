package builtin

import (
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initNull() {
	InitialShape.FuncerStack = InitialShape.FuncerStack.PushAll([]shapes.Funcer{
		shapes.SimpleFuncer(
			types.AnyType,
			"null",
			nil,
			types.NullType,
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				return &values.NullValue{}
			},
		),
	})
}
