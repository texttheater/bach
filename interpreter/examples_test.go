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

func TestArrayTypeExamples(t *testing.T) {
	for _, example := range interpreter.ArrayTypeExamples {
		interpreter.TestProgramStr(
			example.Program,
			example.OutputType,
			example.OutputValue,
			example.Error,
			t,
		)
	}
}

func TestObjectTypeExamples(t *testing.T) {
	for _, example := range interpreter.ObjectTypeExamples {
		interpreter.TestProgramStr(
			example.Program,
			example.OutputType,
			example.OutputValue,
			example.Error,
			t,
		)
	}
}

func TestUnionTypeExamples(t *testing.T) {
	for _, example := range interpreter.UnionTypeExamples {
		interpreter.TestProgramStr(
			example.Program,
			example.OutputType,
			example.OutputValue,
			example.Error,
			t,
		)
	}
}
