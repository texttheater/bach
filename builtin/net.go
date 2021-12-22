package builtin

import (
	"net/url"

	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initNet() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.Funcer{
		expressions.SimpleFuncer(
			types.Str{},
			"urlPathEscape",
			nil,
			types.Str{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				str := string(inputValue.(states.StrValue))
				return states.StrValue(url.PathEscape(str)), nil
			},
		),
		expressions.SimpleFuncer(
			types.Str{},
			"urlPathUnescape",
			nil,
			types.Str{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				str := string(inputValue.(states.StrValue))
				v, err := url.PathUnescape(str)
				if err != nil {
					return nil, err
				}
				return states.StrValue(v), nil
			},
		),
		expressions.SimpleFuncer(
			types.Str{},
			"urlQueryEscape",
			nil,
			types.Str{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				str := string(inputValue.(states.StrValue))
				return states.StrValue(url.QueryEscape(str)), nil
			},
		),
		expressions.SimpleFuncer(
			types.Str{},
			"urlQueryUnescape",
			nil,
			types.Str{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				str := string(inputValue.(states.StrValue))
				v, err := url.QueryUnescape(str)
				if err != nil {
					return nil, err
				}
				return states.StrValue(v), nil
			},
		),
	})
}
