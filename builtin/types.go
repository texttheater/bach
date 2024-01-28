package builtin

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var TypeFuncers = []shapes.Funcer{
	shapes.Funcer{InputType: types.NewVar("A", types.Any{}), Name: "type", Params: nil, OutputType: types.Str{}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		return states.ThunkFromValue(states.StrValue(bindings["A"].String()))
	}, IDs: nil},
}
