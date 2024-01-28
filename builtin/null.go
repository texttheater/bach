package builtin

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var NullFuncers = []shapes.Funcer{
	shapes.Funcer{
		InputType:  types.Any{},
		Name:       "null",
		Params:     nil,
		OutputType: types.Null{},
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return states.ThunkFromValue(states.NullValue{})
		},
	},
}
