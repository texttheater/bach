package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/interpreter"
)

func TestSimpleTypeExamples(t *testing.T) {
	for _, example := range interpreter.SimpleTypeExamples {
		interpreter.TestProgramStr(
			example.Program,
			example.OutputType,
			example.OutputValue,
			example.Error,
			t,
		)
	}
}

func TestTupleTypeExamples(t *testing.T) {
	for _, example := range interpreter.TupleTypeExamples {
		interpreter.TestProgramStr(
			example.Program,
			example.OutputType,
			example.OutputValue,
			example.Error,
			t,
		)
	}
}
