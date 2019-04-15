package patterns

import (
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
)

// pattern/matcher kinda analogous to expression/action

type Pattern interface {
	Typecheck(shapes.Shape) (shapes.Shape, Matcher, error)
}

type Matcher func(states.State) (states.State, bool)
