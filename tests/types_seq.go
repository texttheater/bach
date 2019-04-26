package tests

import (
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func SequenceTypeTestCases() []TestCase {
	return []TestCase{
		{`range(0, 5) type`, types.StrType{}, values.StrValue("Seq<Num>"), nil},
	}
}
