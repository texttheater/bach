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

var ArrFuncers = []expressions.Funcer{
	// for Arr<<A>> +(Arr<<B>>) Arr<<A|B>>
	expressions.RegularFuncer(
		types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		"+",
		[]*params.Param{
			params.SimpleParam("other", types.NewArr(
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
	// for Arr<<A>> drop(Num) Arr<<A>>
	expressions.RegularFuncer(
		types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		"drop",
		[]*params.Param{
			params.SimpleParam("n", types.Num{}),
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
	// for Arr<<A>> dropWhile(for <<A>> Obj<yes: B>|Obj<no: C>) Arr<<A>>
	expressions.RegularFuncer(
		types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		"dropWhile",
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
	// for Arr<<A>> each(for <<A>> <<B>>) Arr<<B>>
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
	// for Arr<<A>> enum(Num) Arr<Tup<Num, <A>>>
	expressions.RegularFuncer(
		types.NewArr(types.NewVar("A", types.Any{})),
		"enum",
		[]*params.Param{
			params.SimpleParam("start", types.Num{}),
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
	// for Arr<<A>> enum Arr<Tup<Num, <A>>>
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
	// for Arr<Bool> every Bool
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
	// for Arr<<A>> find(for <A> Obj<yes: <B>|Obj<no: C>) Null|Tup<Num, <B>>
	expressions.RegularFuncer(
		types.NewArr(types.NewVar("A", types.Any{})),
		"find",
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
		types.NewUnion(
			types.Null{},
			types.NewNearr(
				[]types.Type{
					types.Num{},
					types.NewVar("A", types.Any{}),
				},
				types.VoidArr,
			),
		),
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		},
		nil,
	),
	// for Arr<<A>> findLast(for <A> Obj<yes: <B>|Obj<no: C>) Null|Tup<Num, <B>>
	expressions.RegularFuncer(
		types.NewArr(types.NewVar("A", types.Any{})),
		"findLast",
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
		types.NewUnion(
			types.Null{},
			types.NewNearr(
				[]types.Type{
					types.Num{},
					types.NewVar("A", types.Any{}),
				},
				types.VoidArr,
			),
		),
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		},
		nil,
	),
	// for Arr<<A>> fold(<B>, for <B> (<A>) <B>) <B>
	expressions.RegularFuncer(
		types.NewArr(types.NewVar("A", types.Any{})),
		"fold",
		[]*params.Param{
			params.SimpleParam("start", types.NewVar("B", types.Any{})),
			{
				InputType: types.NewVar("B", types.Any{}),
				Params: []*params.Param{
					params.SimpleParam("next", types.NewVar("A", types.Any{})),
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
	// for Arr<<A>> get(Num) <A>
	expressions.RegularFuncer(
		types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		"get",
		[]*params.Param{
			params.SimpleParam("index", types.Num{}),
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
	// for Arr<Arr<<A>>> join Arr<<A>>
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
	// for Arr<<A>> keep(for <A> Obj<yes: <B>>|Obj<no: <C>>) Arr<B>
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
	// for Arr<<A>> len Num
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
	// for Arr<<A>> max(for <A> (<A>) Bool, <A>) <A>
	expressions.RegularFuncer(
		types.NewArr(types.NewVar("A", types.Any{})),
		"max",
		[]*params.Param{
			{
				InputType: types.NewVar("A", types.Any{}),
				Params: []*params.Param{
					params.SimpleParam("other", types.NewVar("A", types.Any{})),
				},
				OutputType: types.Bool{},
			},
			params.SimpleParam("default", types.NewVar("A", types.Any{})),
		},
		types.NewVar("A", types.Any{}),
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return max(inputState, id, args[0], args[1])
		},
		nil,
	),
	// for Arr<<A>> max(for <A> <B>, for <B> (<B>) Bool, <A>) <A>
	expressions.RegularFuncer(
		types.NewArr(types.NewVar("A", types.Any{})),
		"max",
		[]*params.Param{
			{
				InputType:  types.NewVar("A", types.Any{}),
				Params:     nil,
				OutputType: types.NewVar("B", types.Any{}),
			},
			{
				InputType: types.NewVar("B", types.Any{}),
				Params: []*params.Param{
					params.SimpleParam("other", types.NewVar("B", types.Any{})),
				},
				OutputType: types.Bool{},
			},
			params.SimpleParam("default", types.NewVar("A", types.Any{})),
		},
		types.NewVar("A", types.Any{}),
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return max(inputState, args[0], args[1], args[2])
		},
		nil,
	),
	// for Arr<<A>> min(for <A> (<A>) Bool, <A>) <A>
	expressions.RegularFuncer(
		types.NewArr(types.NewVar("A", types.Any{})),
		"min",
		[]*params.Param{
			{
				InputType: types.NewVar("A", types.Any{}),
				Params: []*params.Param{
					params.SimpleParam("other", types.NewVar("A", types.Any{})),
				},
				OutputType: types.Bool{},
			},
			params.SimpleParam("default", types.NewVar("A", types.Any{})),
		},
		types.NewVar("A", types.Any{}),
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return min(inputState, id, args[0], args[1])
		},
		nil,
	),
	// for Arr<<A>> max(for <A> <B>, for <B> (<B>) Bool, <A>) <A>
	expressions.RegularFuncer(
		types.NewArr(types.NewVar("A", types.Any{})),
		"min",
		[]*params.Param{
			{
				InputType:  types.NewVar("A", types.Any{}),
				Params:     nil,
				OutputType: types.NewVar("B", types.Any{}),
			},
			{
				InputType: types.NewVar("B", types.Any{}),
				Params: []*params.Param{
					params.SimpleParam("other", types.NewVar("B", types.Any{})),
				},
				OutputType: types.Bool{},
			},
			params.SimpleParam("default", types.NewVar("A", types.Any{})),
		},
		types.NewVar("A", types.Any{}),
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return min(inputState, args[0], args[1], args[2])
		},
		nil,
	),
	// for Any range(Num, Num) Arr<Num>
	expressions.RegularFuncer(
		types.Any{},
		"range",
		[]*params.Param{
			params.SimpleParam("from", types.Num{}),
			params.SimpleParam("to", types.Num{}),
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
	// for Arr<<A>> repeat(Num) Arr<<A>>
	expressions.RegularFuncer(
		types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		"repeat",
		[]*params.Param{
			params.SimpleParam("times", types.Num{}),
		},
		types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		nil,
	),
	// for Arr<<A>> rev Arr<<A>>
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
	// for Arr<Bool> some Bool
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
	// for Arr<Num> sort Arr<Num>
	expressions.RegularFuncer(
		types.NewArr(types.Num{}),
		"sort",
		nil,
		types.NewArr(types.Num{}),
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return mySort(inputState, id, numLess)
		},
		nil,
	),
	// for Arr<Str> sort Arr<Str>
	expressions.RegularFuncer(
		types.NewArr(types.Str{}),
		"sort",
		nil,
		types.NewArr(types.Str{}),
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return mySort(inputState, id, strLess)
		},
		nil,
	),
	// for Arr<<A>> sort(for <A> (<A>) Bool) Arr<<A>>
	expressions.RegularFuncer(
		types.NewArr(types.NewVar("A", types.Any{})),
		"sort",
		[]*params.Param{
			{
				InputType: types.NewVar("A", types.Any{}),
				Params: []*params.Param{
					params.SimpleParam("other", types.NewVar("A", types.Any{})),
				},
				OutputType: types.Bool{},
			},
		},
		types.NewArr(types.NewVar("A", types.Any{})),
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return mySort(inputState, id, args[0])
		},
		nil,
	),
	// for Arr<<A>> sortByNum(for <A> Num) Arr<<A>>
	expressions.RegularFuncer(
		types.NewArr(types.NewVar("A", types.Any{})),
		"sortByNum",
		[]*params.Param{
			{
				InputType:  types.NewVar("A", types.Any{}),
				Params:     nil,
				OutputType: types.Num{},
			},
		},
		types.NewArr(types.NewVar("A", types.Any{})),
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return mySort(inputState, args[0], numLess)
		},
		nil,
	),
	// for Arr<<A>> sortByStr(for <A> Str) Arr<<A>>
	expressions.RegularFuncer(
		types.NewArr(types.NewVar("A", types.Any{})),
		"sortByStr",
		[]*params.Param{
			{
				InputType:  types.NewVar("A", types.Any{}),
				Params:     nil,
				OutputType: types.Str{},
			},
		},
		types.NewArr(types.NewVar("A", types.Any{})),
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return mySort(inputState, args[0], strLess)
		},
		nil,
	),
	// for Arr<<A>> sortBy(for <A> <B>, for <B> (<B>) Bool) Arr<<A>>
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
					params.SimpleParam("other", types.NewVar("B", types.Any{})),
				},
				OutputType: types.Bool{},
			},
		},
		types.NewArr(types.NewVar("A", types.Any{})),
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return mySort(inputState, args[0], args[1])
		},
		nil,
	),
	// for Arr<<A>> take(Num) Arr<<A>>
	expressions.RegularFuncer(
		types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		"take",
		[]*params.Param{
			params.SimpleParam("n", types.Num{}),
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
	// for Arr<<A>> takeWhile(for <<A>> Obj<yes: B>|Obj<no: C>) Arr<<A>>
	expressions.RegularFuncer(
		types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		"takeWhile",
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
			types.NewVar("B", types.Any{}),
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
