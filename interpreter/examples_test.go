package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/builtin"
	"github.com/texttheater/bach/interpreter"
)

func TestExamples(t *testing.T) {
	stack := builtin.InitialShape.Stack
	for stack != nil {
		for _, example := range stack.Head.Examples {
			interpreter.TestProgramStr(
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
