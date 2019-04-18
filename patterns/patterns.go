package patterns

import (
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

// pattern/matcher kinda analogous to expression/action

type Pattern interface {
	Typecheck(inputShape shapes.Shape) (outputShape shapes.Shape, restType types.Type, matcher Matcher, err error)
}

type Matcher func(states.State) (*states.VariableStack, bool)

var zeroShape shapes.Shape
