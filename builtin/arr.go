package builtin

import (
	"sort"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initArr() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.Funcer{
		expressions.RegularFuncer(
			types.NewArr(
				types.Str{},
			),
			"sort",
			nil,
			types.NewArr(
				types.Str{},
			),
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
			types.NewArr(
				types.Num{},
			),
			"sort",
			nil,
			types.NewArr(
				types.Num{},
			),
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
			types.NewArr(
				types.NewVar("A", types.Any{}),
			),
			"+",
			[]*parameters.Parameter{
				parameters.SimpleParam(types.NewArr(
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
			types.NewArr(
				types.NewVar("A", types.Any{}),
			),
			"drop",
			[]*parameters.Parameter{
				parameters.SimpleParam(types.Num{}),
			},
			types.NewArr(
				types.NewVar("A", types.Any{}),
			),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				from, err := states.NumFromAction(inputState.Clear(), args[0])
				if err != nil {
					return states.ThunkFromError(err)
				}
				intFrom := int(from)
				if float64(intFrom) != from {
					return states.ThunkFromError(errors.ValueError(
						errors.Code(errors.BadIndex),
						errors.GotValue(states.NumValue(from)),
					))
				}
				input := states.IterFromValue(inputState.Value)
				for i := 0; i < intFrom; i++ {
					_, ok, err := input()
					if err != nil {
						return states.ThunkFromError(err)
					}
					if !ok {
						break
					}
				}
				return states.ThunkFromIter(input)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.NewArr(
				types.NewVar("A", types.Any{}),
			),
			"take",
			[]*parameters.Parameter{
				parameters.SimpleParam(types.Num{}),
			},
			types.NewArr(
				types.NewVar("A", types.Any{}),
			),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				to, err := states.NumFromAction(inputState.Clear(), args[0])
				if err != nil {
					return states.ThunkFromError(err)
				}
				intTo := int(to)
				if float64(intTo) != to {
					return states.ThunkFromError(errors.ValueError(
						errors.Code(errors.BadIndex),
						errors.GotValue(states.NumValue(to)),
					))
				}
				input := states.IterFromValue(inputState.Value)
				output := func() (states.Value, bool, error) {
					if intTo <= 0 {
						return nil, false, nil
					}
					el, ok, err := input()
					if err != nil {
						return nil, false, err
					}
					if !ok {
						return nil, false, nil
					}
					intTo--
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
			"slice",
			[]*parameters.Parameter{
				parameters.SimpleParam(types.Num{}),
				parameters.SimpleParam(types.Num{}),
			},
			types.NewArr(
				types.NewVar("A", types.Any{}),
			),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				argInputState := inputState.Clear()
				from, err := states.NumFromAction(argInputState, args[0])
				if err != nil {
					return states.ThunkFromError(err)
				}
				intFrom := int(from)
				if float64(intFrom) != from || intFrom < 0 {
					return states.ThunkFromError(errors.ValueError(
						errors.Code(errors.BadIndex),
						errors.GotValue(states.NumValue(from)),
					))
				}
				to, err := states.NumFromAction(argInputState, args[1])
				if err != nil {
					return states.ThunkFromError(err)
				}
				intTo := int(to)
				if float64(intTo) != to || intTo < 0 {
					return states.ThunkFromError(errors.ValueError(
						errors.Code(errors.BadIndex),
						errors.GotValue(states.NumValue(to)),
					))
				}
				input := states.IterFromValue(inputState.Value)
				// drop
				for i := 0; i < intFrom; i++ {
					_, ok, err := input()
					if err != nil {
						return states.ThunkFromError(err)
					}
					if !ok {
						break
					}
				}
				// take
				output := func() (states.Value, bool, error) {
					if intTo <= 0 {
						return nil, false, nil
					}
					el, ok, err := input()
					if err != nil {
						return nil, false, err
					}
					if !ok {
						return nil, false, nil
					}
					intTo--
					return el, true, nil
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
			[]*parameters.Parameter{
				parameters.SimpleParam(types.Num{}),
				parameters.SimpleParam(types.Num{}),
			},
			types.NewArr(
				types.Num{},
			),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				argInputState := inputState.Clear()
				from, err := states.NumFromAction(argInputState, args[0])
				if err != nil {
					return states.ThunkFromError(err)
				}
				to, err := states.NumFromAction(argInputState, args[1])
				if err != nil {
					return states.ThunkFromError(err)
				}
				i := from
				var iter func() (states.Value, bool, error)
				iter = func() (states.Value, bool, error) {
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
			"get",
			[]*parameters.Parameter{
				parameters.SimpleParam(types.Num{}),
			},
			types.NewVar("A", types.Any{}),
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				index, err := states.NumFromAction(inputState.Clear(), args[0])
				if err != nil {
					return states.ThunkFromError(err)
				}
				intIndex := int(index)
				if float64(intIndex) != index || intIndex < 0 {
					return states.ThunkFromError(errors.ValueError(
						errors.Code(errors.BadIndex),
						errors.GotValue(states.NumValue(index)),
					))
				}
				input := states.IterFromValue(inputState.Value)
				var el states.Value
				for i := 0; i <= intIndex; i++ {
					var ok bool
					var err error
					el, ok, err = input()
					if err != nil {
						return states.ThunkFromError(err)
					}
					if !ok {
						return states.ThunkFromError(errors.ValueError(
							errors.Code(errors.NoSuchIndex),
							errors.Pos(pos),
						))
					}
				}
				return states.ThunkFromValue(el)
			},
			nil,
		),
	})
}
