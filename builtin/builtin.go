package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

var InitialShape functions.Shape = functions.Shape{
	Type: types.NullType{},
}

func init() {
	initNull()
	initIO()
	initLogic()
	initMath()
	initText()
	initSeq()
	initArr()
	initConv()
	initTypes()
	initValues()
}
