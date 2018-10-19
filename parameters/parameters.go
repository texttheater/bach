package parameters

import (
	"github.com/texttheater/bach/types"
)

type Parameter struct {
	InputType  types.Type
	Parameters []*Parameter
	OutputType types.Type
}

func DumbPar(OutputType types.Type) *Parameter {
	return &Parameter{
		&types.AnyType{},
		[]*Parameter{},
		OutputType,
	}
}
