package tests

import (
	"testing"

	"github.com/texttheater/bach/builtin"
)

func TestExamples(t *testing.T) {
	stack := builtin.InitialShape.Stack
	for stack != nil {
		for _, example := range stack.Head.Examples {
			TestProgramStr(
				example.Program,
				example.OutputType,
				example.OutputValue,
				example.Error,
				t,
			)
		}
		stack = stack.Tail
	}
}
