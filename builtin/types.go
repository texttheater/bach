package builtin

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var TypeFuncers = []expressions.Funcer{
	expressions.RegularFuncer(
		types.NewVar("A", types.Any{}),
		"type",
		nil,
		types.Str{},
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return states.ThunkFromValue(states.StrValue(bindings["A"].String()))
		},
		nil,
	),
}
