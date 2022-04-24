package builtin

import (
	"unicode/utf8"

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
			types.NewArr(
				types.NewVar("A", types.Any{}),
			),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				v := inputState.Value.(states.StrValue)
				offset := 0
				iter := func() (states.Value, bool, error) {
					regexpInputState := states.State{
						Value: v,
					}
					val, err := args[0](regexpInputState, nil).Eval()
					if err != nil {
						return nil, false, err
					}
					objValue, ok := (val.(states.ObjValue))
					if !ok {
						return nil, false, nil
					}
					obj := map[string]*states.Thunk(objValue)
					start, _, _ := obj["start"].EvalInt()
					obj["start"] = states.ThunkFromValue(states.NumValue(start + offset))
					group, _, _ := obj["0"].EvalStr()
					length := len(group)
					end := start + length
					offset += end
					v = states.StrValue(string(v)[end:])
					return objValue, true, nil
				}
				return states.ThunkFromIter(iter)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.Str{},
			"split",
			[]*params.Param{
				{
					InputType: types.Str{},
					OutputType: types.NewUnion(
						types.Null{},
						types.Obj{
							Props: map[string]types.Type{
								"start": types.Num{},
								"0":     types.Str{},
							},
							Rest: types.Any{},
						},
					),
				},
			},
			types.NewArr(
				types.Str{},
			),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				return split(inputState, args, bindings, pos)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.Str{},
			"split",
			[]*params.Param{
				{
					InputType: types.Str{},
					OutputType: types.NewUnion(
						types.Null{},
						types.Obj{
							Props: map[string]types.Type{
								"start": types.Num{},
								"0":     types.Str{},
							},
							Rest: types.Any{},
						},
					),
				},
				params.SimpleParam(types.Num{}),
			},
			types.NewArr(
				types.Str{},
			),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				return split(inputState, args, bindings, pos)
			},
			nil,
		),
	})
}

func split(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
	v := inputState.Value.(states.StrValue)
	regexp := args[0]
	maxSplits := -1
	if len(args) > 1 {
		maxSplits, _, _ = args[1](inputState.Clear(), nil).EvalInt()
	}
	splits := 0
	// Edge case 1: if the separator is empty, split the string into
	// codepoints.
	regexpInputState := states.State{
		Value: states.StrValue(""),
	}
	val, err := regexp(regexpInputState, nil).Eval()
	if err != nil {
		return states.ThunkFromError(err)
	}
	_, ok := val.(states.ObjValue)
	if ok {
		iter := func() (states.Value, bool, error) {
			if len(v) == 0 {
				return nil, false, nil
			}
			var l int
			if splits == maxSplits {
				piece := v
				v = ""
				return piece, true, nil
			}
			_, l = utf8.DecodeRuneInString(string(v))
			piece := v[:l]
			v = v[l:]
			splits++
			return piece, true, nil
		}
		return states.ThunkFromIter(iter)
	}
	// Edge case 2: if the string is empty, return a single-element
	// list containing the empty string.
	if v == "" {
		return states.ThunkFromSlice([]states.Value{
			states.StrValue(""),
		})
	}
	// Now for the normal cases.
	sepAtEnd := false
	iter := func() (states.Value, bool, error) {
		if sepAtEnd {
			sepAtEnd = false
			return states.StrValue(""), true, nil
		}
		if len(v) == 0 {
			// end of string reached
			return nil, false, nil
		}
		if splits == maxSplits {
			piece := v
			v = ""
			return piece, true, nil
		}
		regexpInputState := states.State{
			Value: v,
		}
		val, err := regexp(regexpInputState, nil).Eval()
		if err != nil {
			return nil, false, err
		}
		objValue, ok := val.(states.ObjValue)
		if !ok {
			piece := v
			v = ""
			return piece, true, nil
		}
		obj := map[string]*states.Thunk(objValue)
		sepStart, _, _ := obj["start"].EvalInt()
		sep, _, _ := obj["0"].EvalStr()
		sepLength := len(sep)
		sepEnd := sepStart + sepLength
		piece := v[:sepStart]
		v = v[sepEnd:]
		splits++
		if len(v) == 0 && sepLength > 0 {
			sepAtEnd = true
		}
		return piece, true, nil
	}
	return states.ThunkFromIter(iter)
}
