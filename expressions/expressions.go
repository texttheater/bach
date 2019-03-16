package expressions

import (
	"github.com/texttheater/bach/functions"
)

type Expression interface {
	Typecheck(inputShape functions.Shape, params []*functions.Parameter) (outputShape functions.Shape, action functions.Action, err error)
}

var zeroShape functions.Shape
