package builtin

import (
	"github.com/texttheater/bach/values"
)

func Null(inputValue values.Value, argumentValues []values.Value) values.Value {
	return &values.NullValue{}
}
