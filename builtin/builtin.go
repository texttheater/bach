package builtin

import (
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/types"
)

var InitialShape shapes.Shape = shapes.Shape{types.NullType, nil}

func init() {
	initMath()
	initLogic()
	initNull()
	initSeq()
	initIO()
	initConv()
	initTypes()
	initValues()
}
