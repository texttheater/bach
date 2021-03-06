package builtin

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initTypes() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.Funcer{
		expressions.RegularFuncer(
			types.TypeVariable{
				Name:       "A",
				UpperBound: types.AnyType{},
			},
			"type",
			nil,
			types.StrType{},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				return states.ThunkFromValue(states.StrValue(bindings["A"].String()))
			},
			nil,
		),
	})
}
