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
			types.TypeVariable{
				Name:       "A",
				UpperBound: types.AnyType{},
			},
			"id",
			nil,
			types.TypeVariable{
				Name:       "A",
				UpperBound: types.AnyType{},
			},
			func(inputValue states.Value, argValues []states.Value) (states.Value, error) {
				return inputValue, nil
			},
		),
		expressions.SimpleFuncer(
			types.StrType{},
			"parseFloat",
			nil,
			types.NumType{},
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
			types.StrType{},
			"parseInt",
			[]types.Type{
				types.NumType{},
			},
			types.NumType{},
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
			types.StrType{},
			"parseInt",
			nil,
			types.NumType{},
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
