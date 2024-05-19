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
			{`[1, 2, 3] +[4, 5]`, `Arr<Num>`, `[1, 2, 3, 4, 5]`, nil},
			{`[] +[4, 5]`, `Arr<Num>`, `[4, 5]`, nil},
			{`[1, 2, 3] +[]`, `Arr<Num>`, `[1, 2, 3]`, nil},
			{`[] +[]`, `Tup<>`, `[]`, nil},
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
			{`[true] all`, `Bool`, `true`, nil},
			{`[true, true] all`, `Bool`, `true`, nil},
			{`[true, false, true] all`, `Bool`, `false`, nil},
		},
	},
	shapes.Funcer{
		Summary: "Removes the first n elements from the beginning of an array.",
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
			{`[1, 2, 3] drop(-1)`, `Arr<Num>`, `[1, 2, 3]`, nil},
			{`[1, 2, 3] drop(0)`, `Arr<Num>`, `[1, 2, 3]`, nil},
			{`[1, 2, 3] drop(1)`, `Arr<Num>`, `[2, 3]`, nil},
			{`[1, 2, 3] drop(2)`, `Arr<Num>`, `[3]`, nil},
			{`[1, 2, 3] drop(3)`, `Arr<Num>`, `[]`, nil},
			{`[1, 2, 3] drop(4)`, `Arr<Num>`, `[]`, nil},
		},
	},
	shapes.Funcer{
		Summary: "Removes elements satisfying a condition from the beginning of an array.",
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
	shapes.Funcer{
		Summary:          "Applies a function to every element.",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "an array",
		Name:             "each",
		Params: []*params.Param{
			{
				InputType:   types.NewVar("A", types.Any{}),
				Name:        "f",
				Description: "a function to apply to each element of the input",
				Params:      nil,
				OutputType:  types.NewVar("B", types.Any{}),
			},
		},
		OutputType:        types.NewArr(types.NewVar("B", types.Any{})),
		OutputDescription: "a list with the outputs of f applied to each element of the input",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`[1, 2, 3] each(*2)`, `Arr<Num>`, `[2, 4, 6]`, nil},
			{`[{a: 1}, {a: 2}, {b: 3}, {a: 4}] takeWhile(is {a: _}) each(@a)`, `Arr<Num>`, `[1, 2]`, nil},
		},
	},
	shapes.Funcer{
		Summary:          "Pairs each element with a 0-based index.",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "",
		Name:             "enum",
		Params:           nil,
		OutputType: types.NewArr(types.NewTup([]types.Type{
			types.Num{},
			types.NewVar("A", types.Any{}),
		})),
		OutputDescription: "the input with each element replaced with a 2-element array, the second element of which is the original element and the first is its index in the array, counting from 0",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`["a", "b", "c"] enum`, `Arr<Tup<Num, Str>>`, `[[0, "a"], [1, "b"], [2, "c"]]`, nil},
		},
	},
	shapes.Funcer{
		Summary:          "Pairs each element with an index.",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "an array",
		Name:             "enum",
		Params: []*params.Param{
			params.SimpleParam("start", "at which number to start counting", types.Num{}),
		},
		OutputType: types.NewArr(types.NewTup([]types.Type{
			types.Num{},
			types.NewVar("A", types.Any{}),
		})),
		OutputDescription: "the input with each element replaced with a 2-element array, the second element of which is the original element and the first is its index in the array, counting from start",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`["a", "b", "c"] enum(1)`, `Arr<Tup<Num, Str>>`, `[[1, "a"], [2, "b"], [3, "c"]]`, nil},
		},
	},
	shapes.Funcer{
		Summary:          "Finds the index and first element satisfying a condition.",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "an array",
		Name:             "findFirst",
		Params: []*params.Param{
			{
				InputType:   types.NewVar("A", types.Any{}),
				Name:        "test",
				Description: "a test to apply to elements of the input",
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
		OutputType: types.NewUnion(
			types.Null{},
			types.NewNearr(
				[]types.Type{
					types.Num{},
					types.NewVar("A", types.Any{}),
				},
				types.VoidArr,
			),
		),
		OutputDescription: "the first element of the input passing the test, paired with its index, or Null if none",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			idx, val, err := findFirst(inputState, args[0])
			if err != nil {
				return states.ThunkFromError(err)
			}
			if idx == -1 {
				return states.ThunkFromValue(states.NullValue{})
			}
			return states.ThunkFromSlice([]states.Value{
				states.NumValue(idx),
				val,
			})
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`[1, 2, 3] findFirst(is Num with %2 ==0)`, `Null|Tup<Num, Num>`, `[1, 2]`, nil},
			{`[1, 2, 3] findFirst(is Num with %4 ==0)`, `Null|Tup<Num, Num>`, `null`, nil},
		},
	},
	shapes.Funcer{
		Summary:          "Finds the index of the first element satisfying a condition.",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "an array",
		Name:             "findFirstIndex",
		Params: []*params.Param{
			{
				InputType:   types.NewVar("A", types.Any{}),
				Name:        "test",
				Description: "a test to apply to elements of the input",
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
		OutputType: types.NewUnion(
			types.Null{},
			types.Num{},
		),
		OutputDescription: "the index of the first element of the input passing the test, or Null if none",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			idx, _, err := findFirst(inputState, args[0])
			if err != nil {
				return states.ThunkFromError(err)
			}
			if idx == -1 {
				return states.ThunkFromValue(states.NullValue{})
			}
			return states.ThunkFromValue(states.NumValue(idx))
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`[1, 2, 3] findFirstIndex(is Num with %2 ==0)`, `Null|Num`, `1`, nil},
			{`[1, 2, 3] findFirstIndex(is Num with %4 ==0)`, `Null|Num`, `null`, nil},
		},
	},
	shapes.Funcer{
		Summary:          "Finds the first element satisfying a condition.",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "an array",
		Name:             "findFirstValue",
		Params: []*params.Param{
			{
				InputType:   types.NewVar("A", types.Any{}),
				Name:        "test",
				Description: "a test to apply to elements of the input",
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
		OutputType: types.NewUnion(
			types.Null{},
			types.NewVar("A", types.Any{}),
		),
		OutputDescription: "the first element of the input passing the test, or Null if none",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			idx, val, err := findFirst(inputState, args[0])
			if err != nil {
				return states.ThunkFromError(err)
			}
			if idx == -1 {
				return states.ThunkFromValue(states.NullValue{})
			}
			return states.ThunkFromValue(val)
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`[1, 2, 3] findFirstValue(is Num with %2 ==0)`, `Null|Num`, `2`, nil},
			{`[1, 2, 3] findFirstValue(is Num with %4 ==0)`, `Null|Num`, `null`, nil},
		},
	},
	shapes.Funcer{
		Summary:          "Finds the index and last element satisfying a condition.",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "an array",
		Name:             "findLast",
		Params: []*params.Param{
			{
				InputType:   types.NewVar("A", types.Any{}),
				Name:        "test",
				Description: "a test to apply to elements of the input",
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
		OutputType: types.NewUnion(
			types.Null{},
			types.NewNearr(
				[]types.Type{
					types.Num{},
					types.NewVar("A", types.Any{}),
				},
				types.VoidArr,
			),
		),
		OutputDescription: "the last element of the input passing the test, paired with its index, or Null if none",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			idx, val, err := findLast(inputState, args[0])
			if err != nil {
				return states.ThunkFromError(err)
			}
			if idx == -1 {
				return states.ThunkFromValue(states.NullValue{})
			}
			return states.ThunkFromSlice([]states.Value{
				states.NumValue(idx),
				val,
			})
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`[1, 2, 3, 4] findLast(is Num with %2 ==0)`, `Null|Tup<Num, Num>`, `[3, 4]`, nil},
			{`[1, 2, 3, 4] findLast(is Num with %8 ==0)`, `Null|Tup<Num, Num>`, `null`, nil},
		},
	},
	shapes.Funcer{
		Summary:          "Finds the index of the last element satisfying a condition.",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "an array",
		Name:             "findLastIndex",
		Params: []*params.Param{
			{
				InputType:   types.NewVar("A", types.Any{}),
				Name:        "test",
				Description: "a test to apply to elements of the input",
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
		OutputType: types.NewUnion(
			types.Null{},
			types.Num{},
		),
		OutputDescription: "the index of the last element of the input passing the test, or Null if none",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			idx, _, err := findLast(inputState, args[0])
			if err != nil {
				return states.ThunkFromError(err)
			}
			if idx == -1 {
				return states.ThunkFromValue(states.NullValue{})
			}
			return states.ThunkFromValue(states.NumValue(idx))
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`[1, 2, 3, 4] findLastIndex(is Num with %2 ==0)`, `Null|Num`, `3`, nil},
			{`[1, 2, 3, 4] findLastIndex(is Num with %8 ==0)`, `Null|Num`, `null`, nil},
		},
	},
	shapes.Funcer{
		Summary:          "Finds the last element satisfying a condition.",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "an array",
		Name:             "findLastValue",
		Params: []*params.Param{
			{
				InputType:   types.NewVar("A", types.Any{}),
				Name:        "test",
				Description: "a test to apply to elements of the input",
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
		OutputType: types.NewUnion(
			types.Null{},
			types.NewVar("A", types.Any{}),
		),
		OutputDescription: "the last element of the input passing the test, or Null if none",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			idx, val, err := findLast(inputState, args[0])
			if err != nil {
				return states.ThunkFromError(err)
			}
			if idx == -1 {
				return states.ThunkFromValue(states.NullValue{})
			}
			return states.ThunkFromValue(val)
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`[1, 2, 3, 4] findLastValue(is Num with %2 ==0)`, `Null|Num`, `4`, nil},
			{`[1, 2, 3, 4] findLastValue(is Num with %8 ==0)`, `Null|Num`, `null`, nil},
		},
	},
	shapes.Funcer{
		Summary:          "Aggregates an array recursively.",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "an array",
		Name:             "fold",
		Params: []*params.Param{
			params.SimpleParam("start", "initial accumulator", types.NewVar("B", types.Any{})),
			{
				InputType:   types.NewVar("B", types.Any{}),
				Name:        "combine",
				Description: "a function that combines the current accumulator with the next element to produce a new accumulator",
				Params: []*params.Param{
					params.SimpleParam("next", "the next element", types.NewVar("A", types.Any{})),
				},
				OutputType: types.NewVar("B", types.Any{}),
			},
		},
		OutputType:        types.NewVar("B", types.Any{}),
		OutputDescription: "the accumulator after processing all elements",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`[1, 2, 3] fold(0, +)`, `Num`, `6`, nil},
			{`[2, 3, 4] fold(1, *)`, `Num`, `24`, nil},
		},
	},
	shapes.Funcer{
		Summary: "Returns the element at a given index.",
		InputType: types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		InputDescription: "an array",
		Name:             "get",
		Params: []*params.Param{
			params.SimpleParam("index", "a 0-based index into the input", types.Num{}),
		},
		OutputType:        types.NewVar("A", types.Any{}),
		OutputDescription: "the element at the given index",
		Notes:             "Throws BadIndex if the index is invalid or NoSuchIndex if there is not element at the give index.",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`[] get(0)`, `Void`, ``, errors.TypeError(
				errors.Code(errors.VoidProgram),
			)},
			{`["a", "b", "c"] get(0)`, `Str`, `"a"`, nil},
			{`["a", "b", "c"] get(-1)`, `Str`, ``, errors.ValueError(
				errors.Code(errors.BadIndex),
				errors.GotValue(states.NumValue(-1)),
			)},
			{`["a", "b", "c"] get(2)`, `Str`, `"c"`, nil},
			{`["a", "b", "c"] get(3)`, `Str`, ``, errors.ValueError(
				errors.Code(errors.NoSuchIndex),
				errors.GotValue(states.NumValue(3)),
			)},
		},
	},
	shapes.Funcer{
		Summary: "Returns the element at a given index, or a default value.",
		InputType: types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		InputDescription: "an array",
		Name:             "get",
		Params: []*params.Param{
			params.SimpleParam("index", "a 0-based index into the input", types.Num{}),
			params.SimpleParam("default", "default value to retur if there is no element at the given index", types.NewVar("B", types.Any{})),
		},
		OutputType: types.NewUnion(
			types.NewVar("A", types.Any{}),
			types.NewVar("B", types.Any{}),
		),
		OutputDescription: "the element at the given index, or default if there isn't one.",
		Notes:             "Throws BadIndex if the index is invalid.",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
					return args[1](inputState.Clear(), nil)
				}
				arr, err = arr.Tail.EvalArr()
				if err != nil {
					return states.ThunkFromError(err)
				}
			}
			if arr == nil {
				return args[1](inputState.Clear(), nil)
			}
			return states.ThunkFromValue(arr.Head)
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`[] get(0, null)`, `Null`, `null`, nil},
			{`["a", "b", "c"] get(0, null)`, `Str|Null`, `"a"`, nil},
			{`["a", "b", "c"] get(-1, null)`, `Str|Null`, ``, errors.ValueError(
				errors.Code(errors.BadIndex),
				errors.GotValue(states.NumValue(-1)),
			)},
			{`["a", "b", "c"] get(3, null)`, `Str|Null`, `null`, nil},
		},
	},
	shapes.Funcer{
		Summary: "Join several arrays together into one.",
		InputType: types.NewArr(
			types.NewArr(types.NewVar("A", types.Any{})),
		),
		InputDescription:  "an array of arrays",
		Name:              "join",
		Params:            nil,
		OutputType:        types.NewArr(types.NewVar("A", types.Any{})),
		OutputDescription: "All arrays in the input joined together into one.",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`[] join`, `Arr<<A>>`, `[]`, nil}, // FIXME output type should be Tup<>
			{`[[]] join`, `Tup<>`, `[]`, nil},
			{`[[1]] join`, `Arr<Num>`, `[1]`, nil},
			{`[[1, 2]] join`, `Arr<Num>`, `[1, 2]`, nil},
			{`[[], []] join`, `Tup<>`, `[]`, nil},
			{`[[], [1]] join`, `Arr<Num>`, `[1]`, nil},
			{`[[1], [2, 3]] join`, `Arr<Num>`, `[1, 2, 3]`, nil},
			{`[[1], [2, [3]]] join`, `Arr<Num|Tup<Num>>`, `[1, 2, [3]]`, nil},
		},
	},
	shapes.Funcer{
		Summary:          "Keep only elements satisfying a condition.",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "an array",
		Name:             "keep",
		Params: []*params.Param{
			{
				InputType:   types.NewVar("A", types.Any{}),
				Name:        "test",
				Description: "a test to apply to elements of the input",
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
		OutputType:        types.NewArr(types.NewVar("B", types.Any{})),
		OutputDescription: "The input with all elements not passing the test removed.",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{
				`["a", 1, "b", 2, "c", 3] keep(is Num with %2 >0 elis Str)`,
				`Arr<Num|Str>`,
				`["a", 1, "b", "c", 3]`,
				nil,
			},
			{
				`[{n: 1}, {n: 2}, {n: 3}] keep(is {n: n} with n %2 >0)`,
				`Arr<Obj<n: Num, Void>>`,
				`[{n: 1}, {n: 3}]`,
				nil,
			},
			{
				`[1, 2, 3, 4, 5, 6] keep(if %2 ==0) each(*2)`,
				`Arr<Num>`,
				`[4, 8, 12]`,
				nil,
			},
			{
				`[1, 2, 3, 4, 5, 6] keep(if %2 ==0 not) each(id)`,
				`Arr<Num>`,
				`[1, 3, 5]`,
				nil,
			},
			{
				`[1, 2, 3] keep(if %2 ==0 not) each(+1)`,
				`Arr<Num>`,
				`[2, 4]`,
				nil,
			},
			{
				`[1, 2, 3] keep(if false)`,
				`Arr<Num>`,
				`[]`,
				nil,
			},
			{
				`[{n: 1}, 2, {n: 3}] keep(is {n: n}) each(@n)`,
				`Arr<Num>`,
				`[1, 3]`,
				nil,
			},
		},
	},
	shapes.Funcer{
		Summary:           "Computes the length of an array.",
		InputType:         types.AnyArr,
		InputDescription:  "an array",
		Name:              "len",
		Params:            nil,
		OutputType:        types.Num{},
		OutputDescription: "how many elements the input has",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`[] len`, `Num`, `0`, nil},
			{`["a", 2, []] len`, `Num`, `3`, nil},
		},
	},
	shapes.Funcer{
		Summary:          "",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "",
		Name:             "max",
		Params: []*params.Param{
			{
				InputType: types.NewVar("A", types.Any{}),
				Name:      "less",
				Params: []*params.Param{
					params.SimpleParam("other", "", types.NewVar("A", types.Any{})),
				},
				OutputType: types.Bool{},
			},
			params.SimpleParam("default", "", types.NewVar("A", types.Any{})),
		},
		OutputType:        types.NewVar("A", types.Any{}),
		OutputDescription: "",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return max(inputState, id, args[0], args[1])
		},
		IDs:      nil,
		Examples: []shapes.Example{}},

	shapes.Funcer{
		Summary:          "",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "",
		Name:             "max",
		Params: []*params.Param{
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
		},
		OutputType:        types.NewVar("A", types.Any{}),
		OutputDescription: "",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return max(inputState, args[0], args[1], args[2])
		},
		IDs:      nil,
		Examples: []shapes.Example{}},

	shapes.Funcer{
		Summary:          "",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "",
		Name:             "min",
		Params: []*params.Param{
			{
				InputType: types.NewVar("A", types.Any{}),
				Name:      "less",
				Params: []*params.Param{
					params.SimpleParam("other", "", types.NewVar("A", types.Any{})),
				},
				OutputType: types.Bool{},
			},
			params.SimpleParam("default", "", types.NewVar("A", types.Any{})),
		},
		OutputType:        types.NewVar("A", types.Any{}),
		OutputDescription: "",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return min(inputState, id, args[0], args[1])
		},
		IDs:      nil,
		Examples: []shapes.Example{}},

	shapes.Funcer{
		Summary:          "",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "",
		Name:             "min",
		Params: []*params.Param{
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
		},
		OutputType:        types.NewVar("A", types.Any{}),
		OutputDescription: "",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return min(inputState, args[0], args[1], args[2])
		},
		IDs:      nil,
		Examples: []shapes.Example{}},

	shapes.Funcer{
		Summary:          "",
		InputType:        types.Any{},
		InputDescription: "",
		Name:             "range",
		Params: []*params.Param{
			params.SimpleParam("from", "", types.Num{}),
			params.SimpleParam("to", "", types.Num{}),
		},
		OutputType: types.NewArr(
			types.Num{},
		),
		OutputDescription: "",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs:      nil,
		Examples: []shapes.Example{}},

	shapes.Funcer{
		Summary: "",
		InputType: types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		InputDescription: "",
		Name:             "repeat",
		Params: []*params.Param{
			params.SimpleParam("times", "", types.Num{}),
		},
		OutputType: types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		OutputDescription: "",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		},
		IDs:      nil,
		Examples: []shapes.Example{}},

	shapes.Funcer{
		Summary: "",
		InputType: types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		InputDescription: "",
		Name:             "rev",
		Params:           nil,
		OutputType: types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		OutputDescription: "",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs:      nil,
		Examples: []shapes.Example{}},

	shapes.Funcer{
		Summary:           "",
		InputType:         types.NewArr(types.Bool{}),
		InputDescription:  "",
		Name:              "some",
		Params:            nil,
		OutputType:        types.Bool{},
		OutputDescription: "",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs:      nil,
		Examples: []shapes.Example{}},

	shapes.Funcer{
		Summary:           "",
		InputType:         types.NewArr(types.Num{}),
		InputDescription:  "",
		Name:              "sort",
		Params:            nil,
		OutputType:        types.NewArr(types.Num{}),
		OutputDescription: "",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return mySort(inputState, id, numLess)
		},
		IDs:      nil,
		Examples: []shapes.Example{}},

	shapes.Funcer{
		Summary:           "",
		InputType:         types.NewArr(types.Str{}),
		InputDescription:  "",
		Name:              "sort",
		Params:            nil,
		OutputType:        types.NewArr(types.Str{}),
		OutputDescription: "",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return mySort(inputState, id, strLess)
		},
		IDs:      nil,
		Examples: []shapes.Example{}},

	shapes.Funcer{
		Summary:          "",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "",
		Name:             "sort",
		Params: []*params.Param{
			{
				InputType: types.NewVar("A", types.Any{}),
				Name:      "less",
				Params: []*params.Param{
					params.SimpleParam("other", "", types.NewVar("A", types.Any{})),
				},
				OutputType: types.Bool{},
			},
		},
		OutputType:        types.NewArr(types.NewVar("A", types.Any{})),
		OutputDescription: "",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return mySort(inputState, id, args[0])
		},
		IDs:      nil,
		Examples: []shapes.Example{}},

	shapes.Funcer{
		Summary:          "",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "",
		Name:             "sortByNum",
		Params: []*params.Param{
			{
				InputType:  types.NewVar("A", types.Any{}),
				Name:       "sortKey",
				Params:     nil,
				OutputType: types.Num{},
			},
		},
		OutputType:        types.NewArr(types.NewVar("A", types.Any{})),
		OutputDescription: "",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return mySort(inputState, args[0], numLess)
		},
		IDs:      nil,
		Examples: []shapes.Example{}},

	shapes.Funcer{
		Summary:          "",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "",
		Name:             "sortByStr",
		Params: []*params.Param{
			{
				InputType:  types.NewVar("A", types.Any{}),
				Name:       "sortKey",
				Params:     nil,
				OutputType: types.Str{},
			},
		},
		OutputType:        types.NewArr(types.NewVar("A", types.Any{})),
		OutputDescription: "",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return mySort(inputState, args[0], strLess)
		},
		IDs:      nil,
		Examples: []shapes.Example{}},

	shapes.Funcer{
		Summary:          "",
		InputType:        types.NewArr(types.NewVar("A", types.Any{})),
		InputDescription: "",
		Name:             "sortBy",
		Params: []*params.Param{
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
		},
		OutputType:        types.NewArr(types.NewVar("A", types.Any{})),
		OutputDescription: "",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return mySort(inputState, args[0], args[1])
		},
		IDs:      nil,
		Examples: []shapes.Example{}},

	shapes.Funcer{
		Summary: "",
		InputType: types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		InputDescription: "",
		Name:             "take",
		Params: []*params.Param{
			params.SimpleParam("n", "", types.Num{}),
		},
		OutputType: types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		OutputDescription: "",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs:      nil,
		Examples: []shapes.Example{}},

	shapes.Funcer{
		Summary: "",
		InputType: types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		InputDescription: "",
		Name:             "takeWhile",
		Params: []*params.Param{
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
		},
		OutputType: types.NewArr(
			types.NewVar("B", types.Any{}),
		),
		OutputDescription: "",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs:      nil,
		Examples: []shapes.Example{}},
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

func findFirst(inputState states.State, test states.Action) (index int, value states.Value, err error) {
	arr := inputState.Value.(*states.ArrValue)
	i := 0
	for {
		if arr == nil {
			return -1, nil, nil
		}
		obj, err := test(inputState.Replace(arr.Head), nil).EvalObj()
		if err != nil {
			return -1, nil, err
		}
		_, ok := obj["yes"]
		if ok {
			return i, arr.Head, nil
		}
		i++
		arr, err = arr.Tail.EvalArr()
		if err != nil {
			return -1, nil, err
		}
	}
}

func findLast(inputState states.State, test states.Action) (index int, value states.Value, err error) {
	lastIndex := -1
	var lastValue states.Value = nil
	arr := inputState.Value.(*states.ArrValue)
	i := 0
	for {
		if arr == nil {
			return lastIndex, lastValue, nil
		}
		obj, err := test(inputState.Replace(arr.Head), nil).EvalObj()
		if err != nil {
			return -1, nil, err
		}
		_, ok := obj["yes"]
		if ok {
			lastIndex = i
			lastValue = arr.Head
		}
		i++
		arr, err = arr.Tail.EvalArr()
		if err != nil {
			return -1, nil, err
		}
	}
}
