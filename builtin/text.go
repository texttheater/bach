package builtin

import (
	"strings"

	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func initText() {
	InitialShape.FuncerStack = InitialShape.FuncerStack.PushAll([]shapes.Funcer{
		shapes.SimpleFuncer(
			types.StrType{},
			"split",
			nil,
			&types.ArrType{types.StrType{}},
			func(inputValue values.Value, argumentValues []values.Value) values.Value {
				str, _ := inputValue.(values.StrValue)
				fields := strings.Fields(string(str))
				fieldsValue := make([]values.Value, len(fields))
				for i, field := range fields {
					fieldsValue[i] = values.StrValue(field)
				}
				return values.ArrValue(fieldsValue)
			},
		),
	})
}
