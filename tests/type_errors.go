package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/types"
)

func TestTypeErrors(t *testing.T) {
	TestProgram("-1 *2", nil, nil, errors.E(errors.Code(errors.NoSuchFunction), errors.InputType(types.NullType{}), errors.Name("-"), errors.NumParams(1)), t)
	TestProgram("3 <2 +1", nil, nil, errors.E(errors.Code(errors.NoSuchFunction), errors.InputType(types.BoolType{}), errors.Name("+"), errors.NumParams(1)), t)
	TestProgram("+", nil, nil, errors.E(errors.Code(errors.NoSuchFunction), errors.InputType(types.NullType{}), errors.Name("+"), errors.NumParams(0)), t)
	TestProgram("hurz", nil, nil, errors.E(errors.Code(errors.NoSuchFunction)), t)
}
