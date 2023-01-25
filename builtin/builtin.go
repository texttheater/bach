package builtin

import (
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/types"
)

var InitialShape expressions.Shape = expressions.Shape{
	Type:  types.Null{},
	Stack: nil,
}

func init() {
	initNull()
	initIO()
	initLogic()
	initMath()
	initFmt()
	initText()
	initArr()
	initObj()
	initTypes()
	initValues()
	initRegexp()
	initControl()
}
