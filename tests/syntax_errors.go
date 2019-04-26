package tests

import (
	"github.com/texttheater/bach/errors"
)

func SyntaxErrorTestCases() []TestCase {
	return []TestCase{
		{"&", nil, nil, errors.E(errors.Code(errors.Syntax))},
	}
}
