package builtin

import (
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initArr() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.SimpleFuncer(
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
		functions.RegularFuncer(
			types.AnyType{},
			"range",
			[]*parameters.Parameter{
				parameters.SimpleParam(types.NumType{}),
				parameters.SimpleParam(types.NumType{}),
			},
			&types.ArrType{types.NumType{}},
			func(inputState states.State, args []states.Action) *states.Thunk {
				res0 := args[0](inputState, nil).Eval()
				if res0.Error != nil {
					return states.ThunkFromError(res0.Error)
				}
				start := float64(res0.Value.(states.NumValue))
				res1 := args[1](inputState, nil).Eval()
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
		functions.RegularFuncer(
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
			types.Union(
				types.ObjType{
					PropTypeMap: map[string]types.Type{
						"just": types.TypeVariable{
							Name:       "A",
							UpperBound: types.AnyType{},
						},
					},
					RestType: types.AnyType{},
				},
				types.NullType{},
			),
			func(inputState states.State, args []states.Action) *states.Thunk {
				res0 := args[0](inputState, nil).Eval()
				if res0.Error != nil {
					return states.ThunkFromError(res0.Error)
				}
				index := float64(res0.Value.(states.NumValue))
				intIndex := int(index)
				if intIndex < 0 || index != float64(intIndex) {
					return states.ThunkFromValue(states.NullValue{})
				}
				value := inputState.Value.(*states.ArrValue)
				for i := 0; i < intIndex; i++ {
					if value == nil {
						return states.ThunkFromValue(states.NullValue{})
					}
					tail := value.Tail
					res := tail.Eval()
					if res.Error != nil {
						return states.ThunkFromError(res.Error)
					}
					value = res.Value.(*states.ArrValue)
				}
				if value == nil {
					return states.ThunkFromValue(states.NullValue{})
				}
				return states.ThunkFromValue(states.ObjValue(map[string]*states.Thunk{
					"just": states.ThunkFromValue(value.Head),
				}))
			},
			nil,
		),
	})
}
