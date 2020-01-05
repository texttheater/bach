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
					InputType: types.StrType{},
					OutputType: types.TypeVariable{
						Name: "$",
						UpperBound: types.Union(
							types.NullType{},
							types.NewObjType(map[string]types.Type{
								"start": types.NumType{},
								"0":     types.StrType{},
							}),
						),
					},
				},
			},
			types.TypeVariable{
				Name: "$",
			},
			func(inputState states.State, args []states.Action) *states.Thunk {
				return args[0](inputState, nil)
			},
			nil,
		),
		functions.RegularFuncer(
			types.StrType{},
			"findAll",
			[]*functions.Parameter{
				&functions.Parameter{
					InputType: types.StrType{},
					OutputType: types.TypeVariable{
						Name: "$",
						UpperBound: types.Union(
							types.NullType{},
							types.NewObjType(map[string]types.Type{
								"start": types.NumType{},
								"0":     types.StrType{},
							}),
						),
					},
				},
			},
			&types.ArrType{types.TypeVariable{
				Name: "$",
			}},
			func(inputState states.State, args []states.Action) *states.Thunk {
				offset := 0
				v := inputState.Value.(states.StrValue)
				var iter func() (states.Value, bool, error)
				iter = func() (states.Value, bool, error) {
					regexpInputState := states.State{
						Value: v,
					}
					res := args[0](regexpInputState, nil).Eval()
					if res.Error != nil {
						return nil, false, res.Error
					}
					objValue, ok := (res.Value.(states.ObjValue))
					if !ok {
						return nil, false, nil
					}
					obj := map[string]*states.Thunk(objValue)
					start := int(obj["start"].Eval().Value.(states.NumValue))
					obj["start"] = states.ThunkFromValue(states.NumValue(start + offset))
					length := len(string(obj["0"].Eval().Value.(states.StrValue)))
					end := start + length
					offset += end
					v = states.StrValue(string(v)[end:])
					return objValue, true, nil
				}
				return states.ThunkFromIter(iter)
			},
			nil,
		),
	})
}
