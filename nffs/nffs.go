package nffs

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/builtins"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

type entry struct {
	InputType types.Type
	Name      string
	ArgTypes  []types.Type
	Funcer    func([]functions.Function) functions.Function
}

var entries = []entry{
	entry{
		&types.NumberType{},
		"+",
		[]types.Type{&types.NumberType{}},
		func(argFunctions []functions.Function) functions.Function {
			return builtins.Add{argFunctions[0]}
		},
	},
	entry{
		&types.NumberType{},
		"-",
		[]types.Type{&types.NumberType{}},
		func(argFunctions []functions.Function) functions.Function {
			return builtins.Subtract{argFunctions[0]}
		},
	},
	entry{
		&types.NumberType{},
		"*",
		[]types.Type{&types.NumberType{}},
		func(argFunctions []functions.Function) functions.Function {
			return builtins.Multiply{argFunctions[0]}
		},
	},
	entry{
		&types.NumberType{},
		"/",
		[]types.Type{&types.NumberType{}},
		func(argFunctions []functions.Function) functions.Function {
			return builtins.Divide{argFunctions[0]}
		},
	},
}

func Function(Pos lexer.Position, inputType types.Type, name string, argFunctions []functions.Function) (functions.Function, error) {
	argTypes := make([]types.Type, len(argFunctions))
	for i, f := range argFunctions {
		argTypes[i] = f.Type()
	}
Entries:
	for _, e := range entries {
		if !(e.Name == name) {
			continue
		}
		if len(e.ArgTypes) != len(argFunctions) {
			continue
		}
		if !e.InputType.Subsumes(inputType) {
			continue
		}
		for i, argType := range e.ArgTypes {
			if !argType.Subsumes(argTypes[i]) {
				continue Entries
			}
		}
		return e.Funcer(argFunctions), nil
	}
	return nil, lexer.Errorf(Pos, "no function found: for %v %v(%s)", inputType, name, formatArgTypes(argTypes))
}

func formatArgTypes(argTypes []types.Type) string {
	formatted := make([]string, len(argTypes))
	for i, t := range argTypes {
		formatted[i] = fmt.Sprintf("%v", t)
	}
	return strings.Join(formatted, ", ")
}
