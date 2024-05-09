package builtin

import (
	"sort"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var ArrFuncers = []shapes.Funcer{
	shapes.Funcer{
		Summary: "Concatenates two arrays.",
		InputType: types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		InputDescription: "an array",
		Name:             "+",
		Params: []*params.Param{
			params.SimpleParam("other", "another array", types.NewArr(
				types.NewVar("B", types.Any{}),
			)),
		},
		OutputType: types.NewArr(types.NewUnion(
			types.NewVar("A", types.Any{}),
			types.NewVar("B", types.Any{}),
		)),
		OutputDescription: "the concatenation of both arrays",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`["a", "b"] +["c", "d"]`, `Arr<Str>`, `["a", "b", "c", "d"]`, nil},
			{`["a", "b"] +[1, 2]`, `Arr<Str|Num>`, `["a", "b", 1, 2]`, nil},
		},
	},
	shapes.Funcer{
		Summary:           "Checks if all elements are true.",
		InputType:         types.NewArr(types.Bool{}),
		InputDescription:  "an array of boolean values",
		Name:              "all",
		Params:            nil,
		OutputType:        types.Bool{},
		OutputDescription: "true iff all elements of the input are true",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`[] all`, `Bool`, `true`, nil},
			{`[true, true] all`, `Bool`, `true`, nil},
			{`[true, false, true] all`, `Bool`, `false`, nil},
		},
	},
	shapes.Funcer{
		Summary: "Remove the first n elements from the beginning of an array.",
		InputType: types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		InputDescription: "an array",
		Name:             "drop",
		Params: []*params.Param{
			params.SimpleParam("n", "number of elements to remove", types.Num{}),
		},
		OutputType: types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		OutputDescription: "the input with the first n elements removed",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`[1, 2, 3] drop(0)`, `Arr<Num>`, `[1, 2, 3]`, nil},
			{`[1, 2, 3] drop(1)`, `Arr<Num>`, `[2, 3]`, nil},
			{`[1, 2, 3] drop(2)`, `Arr<Num>`, `[3]`, nil},
			{`[1, 2, 3] drop(3)`, `Arr<Num>`, `[]`, nil},
			{`[1, 2, 3] drop(4)`, `Arr<Num>`, `[]`, nil},
		},
	},
	shapes.Funcer{
		Summary: "Remove elements satisfying a condition from the beginning of an array.",
		InputType: types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		InputDescription: "an array",
		Name:             "dropWhile",
		Params: []*params.Param{
			{
				InputType:   types.NewVar("A", types.Any{}),
				Name:        "test",
				Description: "Return true for elements that should be removed.",
				Params:      nil,
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
		OutputType: types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		OutputDescription: "",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`[{a: 1}, {a: 2}, {b: 3}, {a: 4}] dropWhile(is {a: _})`, `Arr<Obj<b: Num, Void>|Obj<a: Num, Void>>`, `[{b: 3}, {a: 4}]`, nil},
		},
	},

	shapes.Funcer{Summary: "", InputType: types.NewArr(types.NewVar("A", types.Any{})), InputDescription: "", Name: "each", Params: []*params.Param{
		{
			InputType:  types.NewVar("A", types.Any{}),
			Name:       "f",
			Params:     nil,
			OutputType: types.NewVar("B", types.Any{}),
		},
	}, OutputType: types.NewArr(types.NewVar("B", types.Any{})), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(types.NewVar("A", types.Any{})), InputDescription: "", Name: "enum", Params: []*params.Param{
		params.SimpleParam("start", "", types.Num{}),
	}, OutputType: types.NewArr(types.NewTup([]types.Type{
		types.Num{},
		types.NewVar("A", types.Any{}),
	})), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(types.NewVar("A", types.Any{})), InputDescription: "", Name: "enum", Params: nil, OutputType: types.NewArr(types.NewTup([]types.Type{
		types.Num{},
		types.NewVar("A", types.Any{}),
	})), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(types.NewVar("A", types.Any{})), InputDescription: "", Name: "find", Params: []*params.Param{
		{
			InputType: types.NewVar("A", types.Any{}),
			Name:      "test",
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
	}, OutputType: types.NewUnion(
		types.Null{},
		types.NewNearr(
			[]types.Type{
				types.Num{},
				types.NewVar("A", types.Any{}),
			},
			types.VoidArr,
		),
	), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		arr := inputState.Value.(*states.ArrValue)
		i := 0
		for {
			if arr == nil {
				return states.ThunkFromValue(states.NullValue{})
			}
			argInputState := inputState.Replace(arr.Head)
			obj, err := args[0](argInputState, nil).EvalObj()
			if err != nil {
				return states.ThunkFromError(err)
			}
			if thunk, ok := obj["yes"]; ok {
				val, err := thunk.Eval()
				if err != nil {
					return states.ThunkFromError(err)
				}
				return states.ThunkFromSlice([]states.Value{
					states.NumValue(i),
					val,
				})
			}
			i += 1
			arr, err = arr.Tail.EvalArr()
			if err != nil {
				return states.ThunkFromError(err)
			}
		}
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(types.NewVar("A", types.Any{})), InputDescription: "", Name: "findLast", Params: []*params.Param{
		{
			InputType: types.NewVar("A", types.Any{}),
			Name:      "test",
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
	}, OutputType: types.NewUnion(
		types.Null{},
		types.NewNearr(
			[]types.Type{
				types.Num{},
				types.NewVar("A", types.Any{}),
			},
			types.VoidArr,
		),
	), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		arr := inputState.Value.(*states.ArrValue)
		i := -1
		var val states.Value
		for {
			if arr == nil {
				break
			}
			argInputState := inputState.Replace(arr.Head)
			obj, err := args[0](argInputState, nil).EvalObj()
			if err != nil {
				return states.ThunkFromError(err)
			}
			if thunk, ok := obj["yes"]; ok {
				val, err = thunk.Eval()
				if err != nil {
					return states.ThunkFromError(err)
				}
			}
			i += 1
			arr, err = arr.Tail.EvalArr()
			if err != nil {
				return states.ThunkFromError(err)
			}
		}
		if val == nil {
			return states.ThunkFromValue(states.NullValue{})
		}
		return states.ThunkFromSlice([]states.Value{
			states.NumValue(i),
			val,
		})
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(types.NewVar("A", types.Any{})), InputDescription: "", Name: "fold", Params: []*params.Param{
		params.SimpleParam("start", "", types.NewVar("B", types.Any{})),
		{
			InputType: types.NewVar("B", types.Any{}),
			Name:      "combine",
			Params: []*params.Param{
				params.SimpleParam("next", "", types.NewVar("A", types.Any{})),
			},
			OutputType: types.NewVar("B", types.Any{}),
		},
	}, OutputType: types.NewVar("B", types.Any{}), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(
		types.NewVar("A", types.Any{}),
	), InputDescription: "", Name: "get", Params: []*params.Param{
		params.SimpleParam("index", "", types.Num{}),
	}, OutputType: types.NewVar("A", types.Any{}), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(
		types.NewArr(types.NewVar("A", types.Any{})),
	), InputDescription: "", Name: "join", Params: nil, OutputType: types.NewArr(types.NewVar("A", types.Any{})), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(types.NewVar("A", types.Any{})), InputDescription: "", Name: "keep", Params: []*params.Param{
		{
			InputType: types.NewVar("A", types.Any{}),
			Name:      "test",
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
	}, OutputType: types.NewArr(types.NewVar("B", types.Any{})), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.AnyArr, InputDescription: "", Name: "len", Params: nil, OutputType: types.Num{}, OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(types.NewVar("A", types.Any{})), InputDescription: "", Name: "max", Params: []*params.Param{
		{
			InputType: types.NewVar("A", types.Any{}),
			Name:      "less",
			Params: []*params.Param{
				params.SimpleParam("other", "", types.NewVar("A", types.Any{})),
			},
			OutputType: types.Bool{},
		},
		params.SimpleParam("default", "", types.NewVar("A", types.Any{})),
	}, OutputType: types.NewVar("A", types.Any{}), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		return max(inputState, id, args[0], args[1])
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(types.NewVar("A", types.Any{})), InputDescription: "", Name: "max", Params: []*params.Param{
		{
			InputType:  types.NewVar("A", types.Any{}),
			Name:       "sortKey",
			Params:     nil,
			OutputType: types.NewVar("B", types.Any{}),
		},
		{
			InputType: types.NewVar("B", types.Any{}),
			Name:      "less",
			Params: []*params.Param{
				params.SimpleParam("other", "", types.NewVar("B", types.Any{})),
			},
			OutputType: types.Bool{},
		},
		params.SimpleParam("default", "", types.NewVar("A", types.Any{})),
	}, OutputType: types.NewVar("A", types.Any{}), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		return max(inputState, args[0], args[1], args[2])
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(types.NewVar("A", types.Any{})), InputDescription: "", Name: "min", Params: []*params.Param{
		{
			InputType: types.NewVar("A", types.Any{}),
			Name:      "less",
			Params: []*params.Param{
				params.SimpleParam("other", "", types.NewVar("A", types.Any{})),
			},
			OutputType: types.Bool{},
		},
		params.SimpleParam("default", "", types.NewVar("A", types.Any{})),
	}, OutputType: types.NewVar("A", types.Any{}), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		return min(inputState, id, args[0], args[1])
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(types.NewVar("A", types.Any{})), InputDescription: "", Name: "min", Params: []*params.Param{
		{
			InputType:  types.NewVar("A", types.Any{}),
			Name:       "sortKey",
			Params:     nil,
			OutputType: types.NewVar("B", types.Any{}),
		},
		{
			InputType: types.NewVar("B", types.Any{}),
			Name:      "less",
			Params: []*params.Param{
				params.SimpleParam("other", "", types.NewVar("B", types.Any{})),
			},
			OutputType: types.Bool{},
		},
		params.SimpleParam("default", "", types.NewVar("A", types.Any{})),
	}, OutputType: types.NewVar("A", types.Any{}), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		return min(inputState, args[0], args[1], args[2])
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.Any{}, InputDescription: "", Name: "range", Params: []*params.Param{
		params.SimpleParam("from", "", types.Num{}),
		params.SimpleParam("to", "", types.Num{}),
	}, OutputType: types.NewArr(
		types.Num{},
	), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(
		types.NewVar("A", types.Any{}),
	), InputDescription: "", Name: "repeat", Params: []*params.Param{
		params.SimpleParam("times", "", types.Num{}),
	}, OutputType: types.NewArr(
		types.NewVar("A", types.Any{}),
	), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		input := inputState.Value.(*states.ArrValue)
		n, err := args[0](inputState.Clear(), nil).EvalNum()
		nInt := int(n)
		if err != nil {
			return states.ThunkFromError(err)
		}
		c := make(chan *states.Thunk)
		go func() {
			for i := 0; i < nInt; i++ {
				arr := input
				for arr != nil {
					c <- states.ThunkFromValue(arr.Head)
					arr, err = arr.Tail.EvalArr()
					if err != nil {
						c <- states.ThunkFromError(err)
						return
					}
				}
			}
			c <- states.ThunkFromValue(nil)
		}()
		return states.ThunkFromChannel(c)
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(
		types.NewVar("A", types.Any{}),
	), InputDescription: "", Name: "rev", Params: nil, OutputType: types.NewArr(
		types.NewVar("A", types.Any{}),
	), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(types.Bool{}), InputDescription: "", Name: "some", Params: nil, OutputType: types.Bool{}, OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(types.Num{}), InputDescription: "", Name: "sort", Params: nil, OutputType: types.NewArr(types.Num{}), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		return mySort(inputState, id, numLess)
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(types.Str{}), InputDescription: "", Name: "sort", Params: nil, OutputType: types.NewArr(types.Str{}), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		return mySort(inputState, id, strLess)
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(types.NewVar("A", types.Any{})), InputDescription: "", Name: "sort", Params: []*params.Param{
		{
			InputType: types.NewVar("A", types.Any{}),
			Name:      "less",
			Params: []*params.Param{
				params.SimpleParam("other", "", types.NewVar("A", types.Any{})),
			},
			OutputType: types.Bool{},
		},
	}, OutputType: types.NewArr(types.NewVar("A", types.Any{})), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		return mySort(inputState, id, args[0])
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(types.NewVar("A", types.Any{})), InputDescription: "", Name: "sortByNum", Params: []*params.Param{
		{
			InputType:  types.NewVar("A", types.Any{}),
			Name:       "sortKey",
			Params:     nil,
			OutputType: types.Num{},
		},
	}, OutputType: types.NewArr(types.NewVar("A", types.Any{})), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		return mySort(inputState, args[0], numLess)
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(types.NewVar("A", types.Any{})), InputDescription: "", Name: "sortByStr", Params: []*params.Param{
		{
			InputType:  types.NewVar("A", types.Any{}),
			Name:       "sortKey",
			Params:     nil,
			OutputType: types.Str{},
		},
	}, OutputType: types.NewArr(types.NewVar("A", types.Any{})), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		return mySort(inputState, args[0], strLess)
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(types.NewVar("A", types.Any{})), InputDescription: "", Name: "sortBy", Params: []*params.Param{
		{
			InputType:  types.NewVar("A", types.Any{}),
			Name:       "sortKey",
			Params:     nil,
			OutputType: types.NewVar("B", types.Any{}),
		},
		{
			InputType: types.NewVar("B", types.Any{}),
			Name:      "less",
			Params: []*params.Param{
				params.SimpleParam("other", "", types.NewVar("B", types.Any{})),
			},
			OutputType: types.Bool{},
		},
	}, OutputType: types.NewArr(types.NewVar("A", types.Any{})), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		return mySort(inputState, args[0], args[1])
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(
		types.NewVar("A", types.Any{}),
	), InputDescription: "", Name: "take", Params: []*params.Param{
		params.SimpleParam("n", "", types.Num{}),
	}, OutputType: types.NewArr(
		types.NewVar("A", types.Any{}),
	), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
	}, IDs: nil, Examples: []shapes.Example{}},

	shapes.Funcer{Summary: "", InputType: types.NewArr(
		types.NewVar("A", types.Any{}),
	), InputDescription: "", Name: "takeWhile", Params: []*params.Param{
		{
			InputType: types.NewVar("A", types.Any{}),
			Name:      "test",
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
	}, OutputType: types.NewArr(
		types.NewVar("B", types.Any{}),
	), OutputDescription: "", Notes: "", Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
	}, IDs: nil, Examples: []shapes.Example{}},
}

func id(inputState states.State, args []states.Action) *states.Thunk {
	return states.ThunkFromValue(inputState.Value)
}

func numLess(inputState states.State, args []states.Action) *states.Thunk {
	a := float64(inputState.Value.(states.NumValue))
	b, err := args[0](inputState.Clear(), nil).EvalNum()
	if err != nil {
		return states.ThunkFromError(err)
	}
	return states.ThunkFromValue(states.BoolValue(a < b))
}

func strLess(inputState states.State, args []states.Action) *states.Thunk {
	a := string(inputState.Value.(states.StrValue))
	b, err := args[0](inputState.Clear(), nil).EvalStr()
	if err != nil {
		return states.ThunkFromError(err)
	}
	return states.ThunkFromValue(states.BoolValue(a < b))
}

func mySort(inputState states.State, key states.Action, less states.Action) *states.Thunk {
	slice, err := states.SliceFromValue(inputState.Value)
	if err != nil {
		return states.ThunkFromError(err)
	}
	myLess := func(i int, j int) bool {
		aKey, err2 := key(inputState.Replace(slice[i]), nil).Eval()
		if err2 != nil {
			err = err2
			return true
		}
		bKey, err2 := key(inputState.Replace(slice[j]), nil).Eval()
		if err2 != nil {
			err = err2
			return true
		}
		less, err2 := less(inputState.Replace(aKey), []states.Action{states.SimpleAction(bKey)}).EvalBool()
		if err2 != nil {
			err = err2
			return true
		}
		return less
	}
	sort.SliceStable(slice, myLess)
	if err != nil {
		return states.ThunkFromError(err)
	}
	return states.ThunkFromSlice(slice)
}

func max(inputState states.State, key states.Action, less states.Action, def states.Action) *states.Thunk {
	arr := inputState.Value.(*states.ArrValue)
	if arr == nil {
		d, err := def(inputState, nil).Eval()
		if err != nil {
			return states.ThunkFromError(err)
		}
		return states.ThunkFromValue(d)
	}
	record, err := key(inputState.Replace(arr.Head), nil).Eval()
	if err != nil {
		return states.ThunkFromError(err)
	}
	recordHolder := arr.Head
	for {
		arr, err = arr.Tail.EvalArr()
		if err != nil {
			return states.ThunkFromError(err)
		}
		if arr == nil {
			return states.ThunkFromValue(recordHolder)
		}
		val, err := key(inputState.Replace(arr.Head), nil).Eval()
		if err != nil {
			return states.ThunkFromError(err)
		}
		l, err := less(
			inputState.Replace(record),
			[]states.Action{states.SimpleAction(val)},
		).EvalBool()
		if err != nil {
			return states.ThunkFromError(err)
		}
		if l {
			record = val
			recordHolder = arr.Head
		}
	}
}

func min(inputState states.State, key states.Action, less states.Action, def states.Action) *states.Thunk {
	arr := inputState.Value.(*states.ArrValue)
	if arr == nil {
		d, err := def(inputState, nil).Eval()
		if err != nil {
			return states.ThunkFromError(err)
		}
		return states.ThunkFromValue(d)
	}
	record, err := key(inputState.Replace(arr.Head), nil).Eval()
	if err != nil {
		return states.ThunkFromError(err)
	}
	recordHolder := arr.Head
	for {
		arr, err = arr.Tail.EvalArr()
		if err != nil {
			return states.ThunkFromError(err)
		}
		if arr == nil {
			return states.ThunkFromValue(recordHolder)
		}
		val, err := key(inputState.Replace(arr.Head), nil).Eval()
		if err != nil {
			return states.ThunkFromError(err)
		}
		l, err := less(
			inputState.Replace(val),
			[]states.Action{states.SimpleAction(record)},
		).EvalBool()
		if err != nil {
			return states.ThunkFromError(err)
		}
		if l {
			record = val
			recordHolder = arr.Head
		}
	}
}
