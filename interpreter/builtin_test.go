package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/builtin"
	"github.com/texttheater/bach/interpreter"
)

func TestBuiltinArr(t *testing.T) {
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

func TestBuiltinControl(t *testing.T) {
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

func TestBuiltinIO(t *testing.T) {
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

func TestBuiltinLogic(t *testing.T) {
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

func TestBuiltinMath(t *testing.T) {
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

func TestBuiltinNull(t *testing.T) {
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

func TestBuiltinObj(t *testing.T) {
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

func TestBuiltinRegexp(t *testing.T) {
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

func TestBuiltinText(t *testing.T) {
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

func TestBuiltinTypes(t *testing.T) {
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

func TestBuiltinValues(t *testing.T) {
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
