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
				argInputState := states.State{
					Value:     states.NullValue{},
					Stack:     inputState.Stack,
					TypeStack: inputState.TypeStack,
				}
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
				&parameters.Parameter{
					InputType: &types.ArrType{
						ElType: types.TypeVariable{
							Name:       "A",
							UpperBound: types.AnyType{},
						},
					},
					Params:     nil,
					OutputType: types.NumType{},
				},
			},
			types.TypeVariable{
				Name:       "A",
				UpperBound: types.AnyType{},
			},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				res0 := args[0](inputState, nil).Eval()
				if res0.Error != nil {
					return states.ThunkFromError(res0.Error)
				}
				index := float64(res0.Value.(states.NumValue))
				intIndex := int(index)
				if index != float64(intIndex) {
					return states.ThunkFromError(errors.E(
						errors.Code(errors.BadIndex),
					))
				}
				value := inputState.Value.(*states.ArrValue)
				// negative index
				if intIndex < 0 {
					revIndex := -intIndex
					buf := make([]states.Value, revIndex)
					bufIndex := 0
					for true {
						if value == nil {
							if buf[bufIndex] == nil {
								return states.ThunkFromError(errors.E(
									errors.Code(errors.NoSuchIndex),
								))
							}
							return states.ThunkFromValue(buf[bufIndex])
						}
						buf[bufIndex] = value.Head
						bufIndex = (bufIndex + 1) % revIndex
						tail := value.Tail
						res := tail.Eval()
						if res.Error != nil {
							return states.ThunkFromError(res.Error)
						}
						value = res.Value.(*states.ArrValue)
					}
				}
				// nonnegative index
				for i := 0; i < intIndex; i++ {
					if value == nil {
						return states.ThunkFromError(errors.E(
							errors.Code(errors.NoSuchIndex),
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
					return states.ThunkFromError(errors.E(
						errors.Code(errors.NoSuchIndex),
					))
					return states.ThunkFromValue(states.NullValue{})
				}
				return states.ThunkFromValue(value.Head)
			},
			nil,
		),
	})
}
