package expressions

import (
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
)

type Expression interface {
	Typecheck(inputShape shapes.Shape, params []*shapes.Parameter) (outputShape shapes.Shape, action states.Action, err error)
}

var zeroShape shapes.Shape
