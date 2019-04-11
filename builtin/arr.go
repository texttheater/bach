package builtin

import (
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initArr() {
	InitialShape.FuncerStack = InitialShape.FuncerStack.PushAll([]shapes.Funcer{
		shapes.SimpleFuncer(
			types.AnyArrType,
			"length",
			nil,
			types.NumType(),
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				length := len(inputValue.Arr())
				return values.NumValue(float64(length))
			},
		),
	})
}
