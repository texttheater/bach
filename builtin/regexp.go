package builtin

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initRegexp() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.RegularFuncer(
			types.StrType{},
			"findFirst",
			[]*parameters.Parameter{
				&parameters.Parameter{
					InputType: types.StrType{},
					OutputType: types.TypeVariable{
						Name: "A",
						UpperBound: types.Union(
							types.ObjType{
								PropTypeMap: map[string]types.Type{
									"just": types.ObjType{
										PropTypeMap: map[string]types.Type{
											"start": types.NumType{},
											"0":     types.StrType{},
										},
										RestType: types.AnyType{},
									},
								},
								RestType: types.AnyType{},
							},
							types.NullType{},
						),
					},
				},
			},
			types.TypeVariable{
				Name: "A",
			},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				return args[0](inputState, nil)
			},
			nil,
		),
		functions.RegularFuncer(
			types.StrType{},
			"findAll",
			[]*parameters.Parameter{
				&parameters.Parameter{
					InputType: types.StrType{},
					OutputType: types.Union(
						types.ObjType{
							PropTypeMap: map[string]types.Type{
								"just": types.TypeVariable{
									Name: "A",
									UpperBound: types.ObjType{
										PropTypeMap: map[string]types.Type{
											"start": types.NumType{},
											"0":     types.StrType{},
										},
										RestType: types.AnyType{},
									},
								},
							},
							RestType: types.AnyType{},
						},
						types.NullType{},
					),
				},
			},
			&types.ArrType{types.TypeVariable{
				Name: "A",
			}},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
					obj, ok := (res.Value.(states.ObjValue))
					if !ok {
						return nil, false, nil
					}
					res = obj["just"].Eval()
					if res.Error != nil {
						return nil, false, res.Error
					}
					obj = res.Value.(states.ObjValue)
					start := int(obj["start"].Eval().Value.(states.NumValue))
					obj["start"] = states.ThunkFromValue(states.NumValue(start + offset))
					length := len(string(obj["0"].Eval().Value.(states.StrValue)))
					end := start + length
					offset += end
					v = states.StrValue(string(v)[end:])
					return obj, true, nil
				}
				return states.ThunkFromIter(iter)
			},
			nil,
		),
	})
}
