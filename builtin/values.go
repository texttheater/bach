package builtin

import (
	"strconv"

	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initValues() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.Funcer{
		expressions.SimpleFuncer(
			types.Var{
				Name:  "A",
				Bound: types.Any{},
			},
			"id",
			nil,
			types.Var{
				Name:  "A",
				Bound: types.Any{},
			},
			func(inputValue states.Value, argValues []states.Value) (states.Value, error) {
				return inputValue, nil
			},
		),
		expressions.SimpleFuncer(
			types.Str{},
			"parseFloat",
			nil,
			types.Num{},
			func(inputValue states.Value, argValues []states.Value) (states.Value, error) {
				s := string(inputValue.(states.StrValue))
				n, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return nil, err
				}
				return states.NumValue(n), nil
			},
		),
		expressions.SimpleFuncer(
			types.Str{},
			"parseInt",
			[]types.Type{
				types.Num{},
			},
			types.Num{},
			func(inputValue states.Value, argValues []states.Value) (states.Value, error) {
				s := string(inputValue.(states.StrValue))
				b := argValues[0].(states.NumValue)
				n, err := strconv.ParseInt(s, int(b), 64)
				if err != nil {
					return nil, err
				}
				return states.NumValue(n), nil
			},
		),
		expressions.SimpleFuncer(
			types.Str{},
			"parseInt",
			nil,
			types.Num{},
			func(inputValue states.Value, argValues []states.Value) (states.Value, error) {
				s := string(inputValue.(states.StrValue))
				b := 10
				n, err := strconv.ParseInt(s, int(b), 64)
				if err != nil {
					return nil, err
				}
				return states.NumValue(n), nil
			},
		),
	})
}
