package builtin

import (
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initControl() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.RegularFuncer(
			types.Union(types.TypeVariable{"$"}, types.NullType{}),
			"must",
			nil,
			types.TypeVariable{"$"},
			func(inputState states.State, args []states.Action) *states.Thunk {
				inhabits, err := inputState.Value.Inhabits(types.TypeVariable{"$"}, inputState.TypeStack)
				if err != nil {
					return states.ThunkFromError(err)
				}
				if !inhabits {
					return states.ThunkFromError(errors.E(
						errors.Code(errors.UnexpectedValue),
						errors.GotValue(inputState.Value),
					))
				}
				return states.ThunkFromValue(inputState.Value)
			},
		),
	})
}
