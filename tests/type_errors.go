package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestTypeErrors(t *testing.T) {
	TestProgram("3 <2 +1", nil, nil, states.E(states.Code(states.NoSuchFunction), states.InputType(types.BoolType{}), states.Name("+"), states.NumParams(1)), t)
	TestProgram("+", nil, nil, states.E(states.Code(states.NoSuchFunction), states.InputType(types.NullType{}), states.Name("+"), states.NumParams(0)), t)
	TestProgram("hurz", nil, nil, states.E(states.Code(states.NoSuchFunction)), t)
}
