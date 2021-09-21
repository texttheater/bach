package builtin

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initObj() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.Funcer{
		expressions.RegularFuncer(
			types.ObjType{
				PropTypeMap: map[string]types.Type{},
				RestType: types.TypeVariable{
					Name:       "A",
					UpperBound: types.AnyType{},
				},
			},
			"get",
			[]*parameters.Parameter{
				parameters.SimpleParam(types.Union(types.StrType{}, types.NumType{})),
			},
			types.TypeVariable{
				Name:       "A",
				UpperBound: types.AnyType{},
			},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				inputValue := inputState.Value.(states.ObjValue)
				res0 := args[0](inputState.Clear(), nil).Eval()
				if res0.Error != nil {
					return states.ThunkFromError(res0.Error)
				}
				prop, err := res0.Value.Str() // TODO ???
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
				res := thunk.Eval()
				if res.Error != nil {
					return states.ThunkFromError(res.Error)
				}
				return states.ThunkFromValue(res.Value)
			},
			nil,
		),
	})
}
