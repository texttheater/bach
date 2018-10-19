package functions

import (
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

///////////////////////////////////////////////////////////////////////////////

type NFF struct {
	InputType  types.Type // TODO type parameters
	Name       string     // TODO namespaces
	Parameters []*parameters.Parameter
	OutputType types.Type
	Kernel     Kernel
}

///////////////////////////////////////////////////////////////////////////////

type Kernel func(inputState State, argumentValues []values.Value) values.Value

///////////////////////////////////////////////////////////////////////////////
