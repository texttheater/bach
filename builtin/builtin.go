package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

var InitialShape functions.Shape = functions.Shape{types.NullType, nil}

func init() {
	initMath()
	initLogic()
	initNull()
	initSeq()
	initConv()
	initIO()
}
