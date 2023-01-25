package builtin

import (
	"github.com/texttheater/bach/expressions"
)

func initFmt() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.Funcer{})
}
