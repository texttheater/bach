package expressions

import (
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
)

type Expression interface {
	Typecheck(inputShape shapes.Shape, params []*parameters.Parameter) (outputShape shapes.Shape, action states.Action, err error)
}

var zeroShape shapes.Shape
