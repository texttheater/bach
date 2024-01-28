package tests_test

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/tests"
	"github.com/texttheater/bach/types"
)

func TestNull(t *testing.T) {
	tests.TestProgram("1 null", types.Null{}, states.NullValue{}, nil, t)
}
