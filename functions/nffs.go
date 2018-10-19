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
	Funcer     Funcer
}

////////////////////////////////////////////////////////////////////////////////

func DumbFuncer(outputType types.Type, kernel Kernel) Funcer {
	return func(argFunctions []Function) Function {
		return &EvaluatorFunction{
			argFunctions,
			outputType,
			kernel,
		}
	}
}

///////////////////////////////////////////////////////////////////////////////

type Funcer func([]Function) Function

///////////////////////////////////////////////////////////////////////////////

type Kernel func(inputState State, argumentValues []values.Value) values.Value

///////////////////////////////////////////////////////////////////////////////
