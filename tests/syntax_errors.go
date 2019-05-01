package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
)

func TestSyntaxErrors(t *testing.T) {
	TestProgram("&", nil, nil, errors.E(errors.Code(errors.Syntax)), t)
}
