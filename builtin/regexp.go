package builtin

import (
	//	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initRegexp() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.RegularFuncer(
			types.StrType{},
			"findFirst",
			[]*functions.Parameter{
				&functions.Parameter{
					InputType:  types.StrType{},
					OutputType: types.TypeVariable{"$"}, // TODO constrain to Objs of some sort?
				},
			},
			types.TypeVariable{"$"},
			func(inputState states.State, args []states.Action) *states.Thunk {
				return args[0](inputState, nil)
			},
			nil,
		),
	})
}
