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
	return func(argFuncers ...Funcer) Function {
		argFunctions := make([]Function, 0, len(argFuncers))
		for _, argFuncer := range argFuncers {
			argFunctions = append(argFunctions, argFuncer())
		}
		return &EvaluatorFunction{
			argFunctions,
			outputType,
			kernel,
		}
	}
}

func NoArgFuncer(function Function) Funcer {
	return func(...Funcer) Function {
		return function
	}
}

///////////////////////////////////////////////////////////////////////////////

type Funcer func(...Funcer) Function

///////////////////////////////////////////////////////////////////////////////

type Kernel func(inputState State, argumentValues []values.Value) values.Value

///////////////////////////////////////////////////////////////////////////////
