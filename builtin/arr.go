package builtin

import (
	"sort"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initArr() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.Funcer{
		expressions.RegularFuncer(
			types.NewArr(
				types.NewVar("A", types.Any{}),
			),
			"+",
			[]*params.Param{
				params.SimpleParam(types.NewArr(
					types.NewVar("B", types.Any{}),
				)),
			},
			types.NewArr(types.NewUnion(
				types.NewVar("A", types.Any{}),
				types.NewVar("B", types.Any{}),
			)),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				input1 := states.IterFromValue(inputState.Value)
				input2 := states.IterFromAction(inputState.Clear(), args[0])
				output := func() (states.Value, bool, error) {
					el, ok, err := input1()
					if err != nil {
						return nil, false, err
					}
					if ok {
						return el, true, nil
					}
					el, ok, err = input2()
					if err != nil {
						return nil, false, err
					}
					if ok {
						return el, true, nil
					}
					return nil, false, nil
				}
				return states.ThunkFromIter(output)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(
				types.NewVar("A", types.Any{}),
			),
			"drop",
			[]*params.Param{
				params.SimpleParam(types.Num{}),
			},
			types.NewArr(
				types.NewVar("A", types.Any{}),
			),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				arr := inputState.Value.(*states.ArrValue)
				start, err := args[0](inputState.Clear(), nil).EvalInt()
				if err != nil {
					return states.ThunkFromError(err)
				}
				for i := 0; i < start && arr != nil; i++ {
					arr, err = arr.Tail.EvalArr()
					if err != nil {
						return states.ThunkFromError(err)
					}
				}
				return states.ThunkFromValue(arr)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(
				types.NewVar("A", types.Any{}),
			),
			"drop",
			[]*params.Param{
				{
					InputType: types.NewVar("A", types.Any{}),
					Params:    nil,
					OutputType: types.NewUnion(
						types.Obj{
							Props: map[string]types.Type{
								"yes": types.NewVar("B", types.Any{}),
							},
							Rest: types.Any{},
						},
						types.Obj{
							Props: map[string]types.Type{
								"no": types.NewVar("C", types.Any{}),
							},
							Rest: types.Any{},
						},
					),
				},
			},
			types.NewArr(
				types.NewVar("A", types.Any{}),
			),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				arr := inputState.Value.(*states.ArrValue)
				for {
					if arr == nil {
						return states.ThunkFromValue(nil)
					}
					argInputState := inputState.Replace(arr.Head)
					obj, err := args[0](argInputState, nil).EvalObj()
					if err != nil {
						return states.ThunkFromError(err)
					}
					_, ok := obj["yes"]
					if !ok {
						return states.ThunkFromValue(arr)
					}
					arr, err = arr.Tail.EvalArr()
					if err != nil {
						return states.ThunkFromError(err)
					}
				}
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(types.NewVar("A", types.Any{})),
			"each",
			[]*params.Param{
				{
					InputType:  types.NewVar("A", types.Any{}),
					Params:     nil,
					OutputType: types.NewVar("B", types.Any{}),
				},
			},
			types.NewArr(types.NewVar("B", types.Any{})),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				input := states.IterFromValue(inputState.Value)
				output := func() (states.Value, bool, error) {
					val, ok, err := input()
					if err != nil {
						return nil, false, err
					}
					if !ok {
						return nil, false, nil
					}
					argInputState := inputState.Replace(val)
					val, err = args[0](argInputState, nil).Eval()
					if err != nil {
						return nil, false, err
					}
					return val, true, nil
				}
				return states.ThunkFromIter(output)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(types.NewVar("A", types.Any{})),
			"enum",
			[]*params.Param{
				params.SimpleParam(types.Num{}),
			},
			types.NewArr(types.NewTup([]types.Type{
				types.Num{},
				types.NewVar("A", types.Any{}),
			})),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				input := states.IterFromValue(inputState.Value)
				i, err := args[0](inputState.Clear(), nil).EvalNum()
				if err != nil {
					return states.ThunkFromError(err)
				}
				output := func() (states.Value, bool, error) {
					val, ok, err := input()
					if err != nil {
						return nil, false, err
					}
					if !ok {
						return nil, false, nil
					}
					num := states.NumValue(i)
					i++
					return states.NewArrValue([]states.Value{
						num,
						val,
					}), true, nil
				}
				return states.ThunkFromIter(output)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(types.NewVar("A", types.Any{})),
			"enum",
			nil,
			types.NewArr(types.NewTup([]types.Type{
				types.Num{},
				types.NewVar("A", types.Any{}),
			})),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				input := states.IterFromValue(inputState.Value)
				i := 0
				output := func() (states.Value, bool, error) {
					val, ok, err := input()
					if err != nil {
						return nil, false, err
					}
					if !ok {
						return nil, false, nil
					}
					num := states.NumValue(i)
					i++
					return states.NewArrValue([]states.Value{
						num,
						val,
					}), true, nil
				}
				return states.ThunkFromIter(output)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(types.Bool{}),
			"every",
			nil,
			types.Bool{},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				input := states.IterFromValue(inputState.Value)
				for {
					val, ok, err := input()
					if err != nil {
						return states.ThunkFromError(err)
					}
					if !ok {
						return states.ThunkFromValue(states.BoolValue(true))
					}
					if !val.(states.BoolValue) {
						return states.ThunkFromValue(states.BoolValue(false))
					}
				}
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(types.NewVar("A", types.Any{})),
			"fold",
			[]*params.Param{
				params.SimpleParam(types.NewVar("B", types.Any{})),
				{
					InputType: types.NewVar("B", types.Any{}),
					Params: []*params.Param{
						params.SimpleParam(types.NewVar("A", types.Any{})),
					},
					OutputType: types.NewVar("B", types.Any{}),
				},
			},
			types.NewVar("B", types.Any{}),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				acc, err := args[0].Eval(inputState.Clear(), nil)
				if err != nil {
					return states.ThunkFromError(err)
				}
				opInputState := states.State{
					Value: nil,
					Stack: inputState.Stack,
				}
				input := states.IterFromValue(inputState.Value)
				for {
					el, ok, err := input()
					if err != nil {
						return states.ThunkFromError(err)
					}
					if !ok {
						return states.ThunkFromValue(acc)
					}
					opInputState.Value = acc
					acc, err = args[1].Eval(
						opInputState,
						[]states.Action{states.SimpleAction(el)},
					)
					if err != nil {
						return states.ThunkFromError(err)
					}
				}
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(
				types.NewVar("A", types.Any{}),
			),
			"get",
			[]*params.Param{
				params.SimpleParam(types.Num{}),
			},
			types.NewVar("A", types.Any{}),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				arr := inputState.Value.(*states.ArrValue)
				index, err := args[0](inputState.Clear(), nil).EvalInt()
				if err != nil {
					return states.ThunkFromError(err)
				}
				if index < 0 {
					return states.ThunkFromError(errors.ValueError(
						errors.Code(errors.BadIndex),
						errors.GotValue(states.NumValue(index)),
						errors.Pos(pos),
					))
				}
				for i := 0; i < index; i++ {
					if arr == nil {
						return states.ThunkFromError(errors.ValueError(
							errors.Code(errors.NoSuchIndex),
							errors.GotValue(states.NumValue(index)),
							errors.Pos(pos),
						))
					}
					arr, err = arr.Tail.EvalArr()
					if err != nil {
						return states.ThunkFromError(err)
					}
				}
				if arr == nil {
					return states.ThunkFromError(errors.ValueError(
						errors.Code(errors.NoSuchIndex),
						errors.GotValue(states.NumValue(index)),
						errors.Pos(pos),
					))
				}
				return states.ThunkFromValue(arr.Head)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(
				types.NewArr(types.NewVar("A", types.Any{})),
			),
			"join",
			nil,
			types.NewArr(types.NewVar("A", types.Any{})),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				input := states.IterFromValue(inputState.Value)
				var arrIter func() (states.Value, bool, error)
				output := func() (states.Value, bool, error) {
					for {
						if arrIter != nil {
							val, ok, err := arrIter()
							if err != nil {
								return nil, false, nil
							}
							if ok {
								return val, true, nil
							}
						}
						arrVal, ok, err := input()
						if err != nil {
							return nil, false, err
						}
						if !ok {
							return nil, false, nil
						}
						arrIter = states.IterFromValue(arrVal)
					}
				}
				return states.ThunkFromIter(output)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(types.NewVar("A", types.Any{})),
			"keep",
			[]*params.Param{
				{
					InputType: types.NewVar("A", types.Any{}),
					Params:    nil,
					OutputType: types.NewUnion(
						types.Obj{
							Props: map[string]types.Type{
								"yes": types.NewVar("B", types.Any{}),
							},
							Rest: types.Any{},
						},
						types.Obj{
							Props: map[string]types.Type{
								"no": types.NewVar("C", types.Any{}),
							},
							Rest: types.Any{},
						},
					),
				},
			},
			types.NewArr(types.NewVar("B", types.Any{})),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				input := states.IterFromValue(inputState.Value)
				output := func() (states.Value, bool, error) {
					for {
						val, ok, err := input()
						if err != nil {
							return nil, false, err
						}
						if !ok {
							return nil, false, nil
						}
						argInputState := inputState.Replace(val)
						obj, err := args[0](argInputState, nil).EvalObj()
						if err != nil {
							return nil, false, err
						}
						if thunk, ok := obj["yes"]; ok {
							val, err = thunk.Eval()
							if err != nil {
								return nil, false, err
							}
							return val, true, nil
						}
					}
				}
				return states.ThunkFromIter(output)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.AnyArr,
			"len",
			nil,
			types.Num{},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				length := 0
				iter := states.IterFromValue(inputState.Value)
				for {
					_, ok, err := iter()
					if err != nil {
						return states.ThunkFromError(err)
					}
					if !ok {
						return states.ThunkFromValue(states.NumValue(length))
					}
					length += 1
				}
			},
			nil,
		),
		expressions.RegularFuncer(
			types.Any{},
			"range",
			[]*params.Param{
				params.SimpleParam(types.Num{}),
				params.SimpleParam(types.Num{}),
			},
			types.NewArr(
				types.Num{},
			),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				argInputState := inputState.Clear()
				from, err := args[0](argInputState, nil).EvalNum()
				if err != nil {
					return states.ThunkFromError(err)
				}
				to, err := args[1](argInputState, nil).EvalNum()
				if err != nil {
					return states.ThunkFromError(err)
				}
				i := from
				iter := func() (states.Value, bool, error) {
					if i >= to {
						return nil, false, nil
					}
					v := states.NumValue(i)
					i++
					return v, true, nil
				}
				return states.ThunkFromIter(iter)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(
				types.NewVar("A", types.Any{}),
			),
			"rev",
			nil,
			types.NewArr(
				types.NewVar("A", types.Any{}),
			),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				input := states.IterFromValue(inputState.Value)
				var outputArr *states.ArrValue
				for {
					el, ok, err := input()
					if err != nil {
						return states.ThunkFromError(err)
					}
					if !ok {
						break
					}
					outputArr = &states.ArrValue{
						Head: el,
						Tail: states.ThunkFromValue(outputArr),
					}
				}
				return states.ThunkFromValue(outputArr)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(types.Bool{}),
			"some",
			nil,
			types.Bool{},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				input := states.IterFromValue(inputState.Value)
				for {
					val, ok, err := input()
					if err != nil {
						return states.ThunkFromError(err)
					}
					if !ok {
						return states.ThunkFromValue(states.BoolValue(false))
					}
					if val.(states.BoolValue) {
						return states.ThunkFromValue(states.BoolValue(true))
					}
				}
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(types.Num{}),
			"sort",
			nil,
			types.NewArr(types.Num{}),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				slice, err := states.SliceFromValue(inputState.Value)
				if err != nil {
					return states.ThunkFromError(err)
				}
				less := func(i, j int) bool {
					a := slice[i].(states.NumValue)
					b := slice[j].(states.NumValue)
					return a < b
				}
				sort.Slice(slice, less)
				return states.ThunkFromSlice(slice)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(types.Str{}),
			"sort",
			nil,
			types.NewArr(types.Str{}),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				slice, err := states.SliceFromValue(inputState.Value)
				if err != nil {
					return states.ThunkFromError(err)
				}
				less := func(i, j int) bool {
					a := slice[i].(states.StrValue)
					b := slice[j].(states.StrValue)
					return a < b
				}
				sort.Slice(slice, less)
				return states.ThunkFromSlice(slice)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(types.NewVar("A", types.Any{})),
			"sort",
			[]*params.Param{
				{
					InputType: types.NewVar("A", types.Any{}),
					Params: []*params.Param{
						params.SimpleParam(types.NewVar("A", types.Any{})),
					},
					OutputType: types.Bool{},
				},
			},
			types.NewArr(types.NewVar("A", types.Any{})),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				slice, err := states.SliceFromValue(inputState.Value)
				if err != nil {
					return states.ThunkFromError(err)
				}
				less := func(i, j int) bool {
					arg0 := states.State{
						Value: slice[i],
						Stack: inputState.Stack,
					}
					arg1 := states.SimpleAction(slice[j])
					val, err2 := args[0](arg0, []states.Action{arg1}).Eval()
					if err2 != nil {
						err = err2
						return true
					}
					return bool(val.(states.BoolValue))
				}
				sort.SliceStable(slice, less)
				if err != nil {
					return states.ThunkFromError(err)
				}
				return states.ThunkFromSlice(slice)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(types.NewVar("A", types.Any{})),
			"sortBy",
			[]*params.Param{
				{
					InputType:  types.NewVar("A", types.Any{}),
					Params:     nil,
					OutputType: types.NewVar("B", types.Any{}),
				},
				{
					InputType: types.NewVar("B", types.Any{}),
					Params: []*params.Param{
						params.SimpleParam(types.NewVar("B", types.Any{})),
					},
					OutputType: types.Bool{},
				},
			},
			types.NewArr(types.NewVar("A", types.Any{})),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				slice, err := states.SliceFromValue(inputState.Value)
				if err != nil {
					return states.ThunkFromError(err)
				}
				less := func(i, j int) bool {
					key := args[0]
					cmp := args[1]
					keyInputState := states.State{
						Value: slice[i],
						Stack: inputState.Stack,
					}
					a, err2 := key(keyInputState, nil).Eval()
					if err2 != nil {
						err = err2
						return true
					}
					keyInputState.Value = slice[j]
					b, err2 := key(keyInputState, nil).Eval()
					if err2 != nil {
						err = err2
						return true
					}
					arg0 := states.State{
						Value: a,
						Stack: inputState.Stack,
					}
					arg1 := states.SimpleAction(b)
					val, err2 := cmp(arg0, []states.Action{arg1}).Eval()
					if err2 != nil {
						err = err2
						return true
					}
					return bool(val.(states.BoolValue))
				}
				sort.SliceStable(slice, less)
				if err != nil {
					return states.ThunkFromError(err)
				}
				return states.ThunkFromSlice(slice)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(
				types.NewVar("A", types.Any{}),
			),
			"take",
			[]*params.Param{
				params.SimpleParam(types.Num{}),
			},
			types.NewArr(
				types.NewVar("A", types.Any{}),
			),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				arr := inputState.Value.(*states.ArrValue)
				stop, err := args[0](inputState.Clear(), nil).EvalInt()
				if err != nil {
					return states.ThunkFromError(err)
				}
				i := 0
				output := func() (states.Value, bool, error) {
					if i >= stop {
						return nil, false, nil
					}
					if arr == nil {
						return nil, false, nil
					}
					el := arr.Head
					i += 1
					arr, err = arr.Tail.EvalArr()
					if err != nil {
						return nil, false, err
					}
					return el, true, nil
				}
				return states.ThunkFromIter(output)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(
				types.NewVar("A", types.Any{}),
			),
			"take",
			[]*params.Param{
				{
					InputType: types.NewVar("A", types.Any{}),
					Params:    nil,
					OutputType: types.NewUnion(
						types.Obj{
							Props: map[string]types.Type{
								"yes": types.NewVar("B", types.Any{}),
							},
							Rest: types.Any{},
						},
						types.Obj{
							Props: map[string]types.Type{
								"no": types.NewVar("C", types.Any{}),
							},
							Rest: types.Any{},
						},
					),
				},
			},
			types.NewArr(
				types.NewVar("B", types.NewVar("A", types.Any{})),
			),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				input := states.IterFromValue(inputState.Value)
				output := func() (states.Value, bool, error) {
					val, ok, err := input()
					if err != nil {
						return nil, false, err
					}
					if !ok {
						return nil, false, nil
					}
					argInputState := inputState.Replace(val)
					obj, err := args[0](argInputState, nil).EvalObj()
					if err != nil {
						return nil, false, err
					}
					if thunk, ok := obj["yes"]; ok {
						val, err = thunk.Eval()
						if err != nil {
							return nil, false, err
						}
						return val, true, nil
					}
					return nil, false, nil
				}
				return states.ThunkFromIter(output)
			},
			nil,
		),
	})
}
