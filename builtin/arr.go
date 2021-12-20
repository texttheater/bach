package builtin

import (
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
			&types.ArrType{types.TypeVariable{"A", types.AnyType{}}},
			"+",
			[]*parameters.Parameter{
				parameters.SimpleParam(&types.ArrType{
					types.TypeVariable{"B", types.AnyType{}},
				}),
			},
			&types.ArrType{types.Union(
				types.TypeVariable{"A", types.AnyType{}},
				types.TypeVariable{"B", types.AnyType{}},
			)},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				arr1 := inputState.Value.(*states.ArrValue)
				res := args[0](inputState.Clear(), nil).Eval()
				if res.Error != nil {
					return states.ThunkFromError(res.Error)
				}
				arr2 := res.Value.(*states.ArrValue)
				iter := func() (states.Value, bool, error) {
					if arr1 != nil {
						head := arr1.Head
						res := arr1.Tail.Eval()
						if res.Error != nil {
							return nil, false, res.Error
						}
						arr1 = res.Value.(*states.ArrValue)
						return head, true, nil
					}
					if arr2 != nil {
						head := arr2.Head
						res := arr2.Tail.Eval()
						if res.Error != nil {
							return nil, false, res.Error
						}
						arr2 = res.Value.(*states.ArrValue)
						return head, true, nil
					}
					return nil, false, nil
				}
				return states.ThunkFromIter(iter)
			},
			nil,
		),
		expressions.SimpleFuncer(
			&types.ArrType{
				ElType: types.TypeVariable{
					Name:       "A",
					UpperBound: types.AnyType{},
				},
			},
			"rev",
			nil,
			&types.ArrType{
				ElType: types.TypeVariable{
					Name:       "A",
					UpperBound: types.AnyType{},
				},
			},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				inputArr := inputValue.(*states.ArrValue)
				var outputArr *states.ArrValue
				for inputArr != nil {
					outputArr = &states.ArrValue{
						Head: inputArr.Head,
						Tail: states.ThunkFromValue(outputArr),
					}
					res := inputArr.Tail.Eval()
					if res.Error != nil {
						return nil, res.Error
					}
					inputArr = res.Value.(*states.ArrValue)
				}
				return outputArr, nil
			},
		),
		expressions.RegularFuncer(
			&types.ArrType{
				ElType: types.TypeVariable{
					Name:       "A",
					UpperBound: types.AnyType{},
				},
			},
			"drop",
			[]*parameters.Parameter{
				parameters.SimpleParam(types.NumType{}),
			},
			&types.ArrType{
				ElType: types.TypeVariable{
					Name:       "A",
					UpperBound: types.AnyType{},
				},
			},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				n, err := states.IntFromAction(inputState.Clear(), args[0])
				if err != nil {
					return states.ThunkFromError(err)
				}
				input := states.IterFromValue(inputState.Value)
				for i := 0; i < n; i++ {
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
			&types.ArrType{
				ElType: types.TypeVariable{
					Name:       "A",
					UpperBound: types.AnyType{},
				},
			},
			"take",
			[]*parameters.Parameter{
				parameters.SimpleParam(types.NumType{}),
			},
			&types.ArrType{
				ElType: types.TypeVariable{
					Name:       "A",
					UpperBound: types.AnyType{},
				},
			},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				res := args[0](inputState.Clear(), nil).Eval()
				if res.Error != nil {
					return states.ThunkFromError(res.Error)
				}
				n := int(res.Value.(states.NumValue))
				input := states.IterFromValue(inputState.Value)
				output := func() (states.Value, bool, error) {
					if n <= 0 {
						return nil, false, nil
					}
					el, ok, err := input()
					if err != nil {
						return nil, false, err
					}
					if !ok {
						return nil, false, nil
					}
					n--
					return el, true, nil
				}
				return states.ThunkFromIter(output)
			},
			nil,
		),
		expressions.RegularFuncer(
			&types.ArrType{
				ElType: types.TypeVariable{
					Name:       "A",
					UpperBound: types.AnyType{},
				},
			},
			"slice",
			[]*parameters.Parameter{
				parameters.SimpleParam(types.NumType{}),
				parameters.SimpleParam(types.NumType{}),
			},
			&types.ArrType{
				ElType: types.TypeVariable{
					Name:       "A",
					UpperBound: types.AnyType{},
				},
			},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				// validate first argument
				argInputState := inputState.Clear()
				res0 := args[0](argInputState, nil).Eval()
				if res0.Error != nil {
					return states.ThunkFromError(res0.Error)
				}
				from := float64(res0.Value.(states.NumValue))
				intFrom := int(from)
				if from != float64(intFrom) || intFrom < 0 {
					return states.ThunkFromError(errors.ValueError(
						errors.Code(errors.BadIndex),
						errors.Pos(pos),
					))
				}
				// validate second argument
				res1 := args[1](argInputState, nil).Eval()
				if res1.Error != nil {
					return states.ThunkFromError(res1.Error)
				}
				to := float64(res1.Value.(states.NumValue))
				intTo := int(to)
				if to != float64(intTo) || intTo < 0 {
					return states.ThunkFromError(errors.ValueError(
						errors.Code(errors.BadIndex),
						errors.Pos(pos),
					))
				}
				// make iter
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
		expressions.SimpleFuncer(
			types.AnyArrType,
			"len",
			nil,
			types.NumType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				length := 0
				iter := states.IterFromValue(inputValue)
				for {
					_, ok, err := iter()
					if err != nil {
						return nil, err
					}
					if !ok {
						return states.NumValue(float64(length)), nil
					}
					length += 1
				}
			},
		),
		expressions.RegularFuncer(
			types.AnyType{},
			"range",
			[]*parameters.Parameter{
				parameters.SimpleParam(types.NumType{}),
				parameters.SimpleParam(types.NumType{}),
			},
			&types.ArrType{types.NumType{}},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				argInputState := inputState.Clear()
				res0 := args[0](argInputState, nil).Eval()
				if res0.Error != nil {
					return states.ThunkFromError(res0.Error)
				}
				start := float64(res0.Value.(states.NumValue))
				res1 := args[1](argInputState, nil).Eval()
				if res1.Error != nil {
					return states.ThunkFromError(res1.Error)
				}
				end := float64(res1.Value.(states.NumValue))
				i := start
				var iter func() (states.Value, bool, error)
				iter = func() (states.Value, bool, error) {
					if i > end {
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
			&types.ArrType{
				ElType: types.TypeVariable{
					Name:       "A",
					UpperBound: types.AnyType{},
				},
			},
			"get",
			[]*parameters.Parameter{
				parameters.SimpleParam(types.NumType{}),
			},
			types.TypeVariable{
				Name:       "A",
				UpperBound: types.AnyType{},
			},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				res0 := args[0](inputState.Clear(), nil).Eval()
				if res0.Error != nil {
					return states.ThunkFromError(res0.Error)
				}
				index := float64(res0.Value.(states.NumValue))
				intIndex := int(index)
				if index != float64(intIndex) || intIndex < 0 {
					return states.ThunkFromError(errors.ValueError(
						errors.Code(errors.BadIndex),
						errors.Pos(pos),
					))
				}
				value := inputState.Value.(*states.ArrValue)
				for i := 0; i < intIndex; i++ {
					if value == nil {
						return states.ThunkFromError(errors.ValueError(
							errors.Code(errors.NoSuchIndex),
							errors.Pos(pos),
						))
					}
					tail := value.Tail
					res := tail.Eval()
					if res.Error != nil {
						return states.ThunkFromError(res.Error)
					}
					value = res.Value.(*states.ArrValue)
				}
				if value == nil {
					return states.ThunkFromError(errors.ValueError(
						errors.Code(errors.NoSuchIndex),
						errors.Pos(pos),
					))
					return states.ThunkFromValue(states.NullValue{})
				}
				return states.ThunkFromValue(value.Head)
			},
			nil,
		),
	})
}
