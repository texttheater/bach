package builtin

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var NullFuncers = []shapes.Funcer{
	shapes.Funcer{
		Summary:           "Returns the null value.",
		InputType:         types.AnyType{},
		InputDescription:  "any value (is ignored)",
		Name:              "null",
		Params:            nil,
		OutputType:        types.NullType{},
		OutputDescription: "the null value (the only value of this type)",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return states.ThunkFromValue(states.NullValue{})
		},
		Examples: []shapes.Example{
			{"null", "Null", "null", nil},
			{"1 null", "Null", "null", nil},
		},
	},
}
