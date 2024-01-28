package tests_test

import (
	"testing"

	"github.com/texttheater/bach/builtin"
	"github.com/texttheater/bach/tests"
)

func TestExamples(t *testing.T) {
	stack := builtin.InitialShape.Stack
	for stack != nil {
		for _, example := range stack.Head.Examples {
			tests.TestProgramStr(
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
