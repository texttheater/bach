package tests

import (
	"testing"

	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestSequenceTypes(t *testing.T) {
	TestProgram(`range(0, 5) type`, types.StrType{}, values.StrValue("Seq<Num>"), nil, t)
}
