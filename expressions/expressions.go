package expressions

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/states"
)

type Expression interface {
	Typecheck(inputShape functions.Shape, params []*parameters.Parameter) (outputShape functions.Shape, action states.Action, err error)
}

var zeroShape functions.Shape
