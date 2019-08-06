package functions

import (
	"github.com/texttheater/bach/states"
)

type Expression interface {
	Typecheck(inputShape Shape, params []*Parameter) (outputShape Shape, action states.Action, err error)
}
