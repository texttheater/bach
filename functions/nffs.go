package functions

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

///////////////////////////////////////////////////////////////////////////////

type NFF struct {
	InputType  types.Type // TODO type parameters
	Name       string     // TODO namespaces
	ArgTypes   []types.Type
	OutputType types.Type
	Kernel     Kernel
}

func (nff NFF) Function(inputShape Shape, name string, argFunctions []Function) (Function, bool) {
	if name != nff.Name {
		return nil, false
	}
	if len(argFunctions) != len(nff.ArgTypes) {
		return nil, false
	}
	if !nff.InputType.Subsumes(inputShape.Type) {
		return nil, false
	}
	for i, argType := range nff.ArgTypes {
		if !argType.Subsumes(argFunctions[i].OutputShape(inputShape.Stack).Type) {
			return nil, false
		}
	}
	return &EvaluatorFunction{
		argFunctions,
		nff.OutputType,
		nff.Kernel,
	}, true
}

///////////////////////////////////////////////////////////////////////////////

type Kernel func(inputState State, argumentValues []values.Value) values.Value

///////////////////////////////////////////////////////////////////////////////
