package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestValues(t *testing.T) {
	TestProgram(
		`false if id then 1 else reject ok`,
		types.NumType{},
		nil,
		errors.E(
			errors.Code(errors.UnexpectedValue),
			errors.GotValue(values.BoolValue(false)),
		),
		t,
	)
}
