package tests

import (
	"testing"

	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestNull(t *testing.T) {
	TestProgram("1 null", types.NullType{}, &values.NullValue{}, nil, t)
}
