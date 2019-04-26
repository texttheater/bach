package tests

import (
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/types"
)

func TypeErrorTestCases() []TestCase {
	return []TestCase{
		{"-1 *2", nil, nil, errors.E(errors.Code(errors.NoSuchFunction), errors.InputType(types.NullType{}), errors.Name("-"), errors.NumParams(1))},
		{"3 <2 +1", nil, nil, errors.E(errors.Code(errors.NoSuchFunction), errors.InputType(types.BoolType{}), errors.Name("+"), errors.NumParams(1))},
		{"+", nil, nil, errors.E(errors.Code(errors.NoSuchFunction), errors.InputType(types.NullType{}), errors.Name("+"), errors.NumParams(0))},
		{"hurz", nil, nil, errors.E(errors.Code(errors.NoSuchFunction))},
	}
}
