package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/builtin"
	"github.com/texttheater/bach/interpreter"
)

func TestExamplesArr(t *testing.T) {
	for _, funcer := range builtin.ArrFuncers {
		for _, example := range funcer.Examples {
			interpreter.TestProgramStr(
				example.Program,
				example.OutputType,
				example.OutputValue,
				example.Error,
				t,
			)
		}
	}
}

func TestExamplesControl(t *testing.T) {
	for _, funcer := range builtin.ControlFuncers {
		for _, example := range funcer.Examples {
			interpreter.TestProgramStr(
				example.Program,
				example.OutputType,
				example.OutputValue,
				example.Error,
				t,
			)
		}
	}
}

func TestExamplesIO(t *testing.T) {
	for _, funcer := range builtin.IOFuncers {
		for _, example := range funcer.Examples {
			interpreter.TestProgramStr(
				example.Program,
				example.OutputType,
				example.OutputValue,
				example.Error,
				t,
			)
		}
	}
}

func TestExamplesLogic(t *testing.T) {
	for _, funcer := range builtin.LogicFuncers {
		for _, example := range funcer.Examples {
			interpreter.TestProgramStr(
				example.Program,
				example.OutputType,
				example.OutputValue,
				example.Error,
				t,
			)
		}
	}
}

func TestExamplesMath(t *testing.T) {
	for _, funcer := range builtin.MathFuncers {
		for _, example := range funcer.Examples {
			interpreter.TestProgramStr(
				example.Program,
				example.OutputType,
				example.OutputValue,
				example.Error,
				t,
			)
		}
	}
}

func TestExamplesNull(t *testing.T) {
	for _, funcer := range builtin.NullFuncers {
		for _, example := range funcer.Examples {
			interpreter.TestProgramStr(
				example.Program,
				example.OutputType,
				example.OutputValue,
				example.Error,
				t,
			)
		}
	}
}

func TestExamplesObj(t *testing.T) {
	for _, funcer := range builtin.ObjFuncers {
		for _, example := range funcer.Examples {
			interpreter.TestProgramStr(
				example.Program,
				example.OutputType,
				example.OutputValue,
				example.Error,
				t,
			)
		}
	}
}

func TestExamplesRegexp(t *testing.T) {
	for _, funcer := range builtin.RegexpFuncers {
		for _, example := range funcer.Examples {
			interpreter.TestProgramStr(
				example.Program,
				example.OutputType,
				example.OutputValue,
				example.Error,
				t,
			)
		}
	}
}

func TestExamplesText(t *testing.T) {
	for _, funcer := range builtin.TextFuncers {
		for _, example := range funcer.Examples {
			interpreter.TestProgramStr(
				example.Program,
				example.OutputType,
				example.OutputValue,
				example.Error,
				t,
			)
		}
	}
}

func TestExamplesTypes(t *testing.T) {
	for _, funcer := range builtin.TypeFuncers {
		for _, example := range funcer.Examples {
			interpreter.TestProgramStr(
				example.Program,
				example.OutputType,
				example.OutputValue,
				example.Error,
				t,
			)
		}
	}
}

func TestExamplesValues(t *testing.T) {
	for _, funcer := range builtin.ValueFuncers {
		for _, example := range funcer.Examples {
			interpreter.TestProgramStr(
				example.Program,
				example.OutputType,
				example.OutputValue,
				example.Error,
				t,
			)
		}
	}
}
