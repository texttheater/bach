package builtin

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initObj() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.Funcer{
		expressions.RegularFuncer(
			types.Obj{
				Props: map[string]types.Type{},
				Rest:  types.NewVar("A", types.Any{}),
			},
			"get",
			[]*params.Param{
				params.SimpleParam(types.NewUnion(types.Str{}, types.Num{})),
			},
			types.NewVar("A", types.Any{}),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				inputValue := inputState.Value.(states.ObjValue)
				val, err := args[0](inputState.Clear(), nil).Eval()
				if err != nil {
					return states.ThunkFromError(err)
				}
				prop, err := val.Str() // TODO ???
				if err != nil {
					return states.ThunkFromError(err)
				}
				thunk, ok := inputValue[prop]
				if !ok {
					return states.ThunkFromError(errors.ValueError(
						errors.Code(errors.NoSuchProperty),
						errors.Pos(pos),
					))
				}
				val, err = thunk.Eval()
				if err != nil {
					return states.ThunkFromError(err)
				}
				return states.ThunkFromValue(val)
			},
			nil,
		),
	})
}
