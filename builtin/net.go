package builtin

import (
	"net/url"

	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initNet() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]functions.Funcer{
		functions.SimpleFuncer(
			types.StrType{},
			"urlPathEscape",
			nil,
			types.StrType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				str := string(inputValue.(states.StrValue))
				return states.StrValue(url.PathEscape(str)), nil
			},
		),
		functions.SimpleFuncer(
			types.StrType{},
			"urlPathUnescape",
			nil,
			types.Union(types.StrType{}, types.NullType{}),
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				str := string(inputValue.(states.StrValue))
				v, err := url.PathUnescape(str)
				if err != nil {
					return states.NullValue{}, nil
				}
				return states.StrValue(v), nil
			},
		),
		functions.SimpleFuncer(
			types.StrType{},
			"urlQueryEscape",
			nil,
			types.StrType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				str := string(inputValue.(states.StrValue))
				return states.StrValue(url.QueryEscape(str)), nil
			},
		),
		functions.SimpleFuncer(
			types.StrType{},
			"urlQueryUnescape",
			nil,
			types.StrType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				str := string(inputValue.(states.StrValue))
				v, err := url.QueryUnescape(str)
				if err != nil {
					return states.NullValue{}, nil
				}
				return states.StrValue(v), nil
			},
		),
	})
}
