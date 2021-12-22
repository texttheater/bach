package tests

import (
	"testing"

	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestNull(t *testing.T) {
	TestProgram("1 null", types.Null{}, states.NullValue{}, nil, t)
}
