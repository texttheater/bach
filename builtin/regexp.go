package builtin

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initRegexp() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.Funcer{
		expressions.RegularFuncer(
			types.Str{},
			"findFirst",
			[]*params.Param{
				{
					InputType: types.Str{},
					OutputType: types.NewVar("A", types.NewUnion(
						types.Null{},
						types.Obj{
							Props: map[string]types.Type{
								"start": types.Num{},
								"0":     types.Str{},
							},
							Rest: types.Any{},
						},
					)),
				},
			},
			types.Var{
				Name: "A",
			},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				return args[0](inputState, nil)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.Str{},
			"findAll",
			[]*params.Param{
				{
					InputType: types.Str{},
					OutputType: types.NewVar("A", types.NewUnion(
						types.Null{},
						types.Obj{
							Props: map[string]types.Type{
								"start": types.Num{},
								"0":     types.Str{},
							},
							Rest: types.Any{},
						},
					)),
				},
			},
			types.NewArr(types.Var{
				Name: "A",
			}),

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
