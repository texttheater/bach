package builtin

import (
	"bytes"
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initText() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.Funcer{
		expressions.RegularFuncer(
			types.StrType{},
			"bytes",
			nil,
			&types.ArrType{types.NumType{}},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				str := string(inputState.Value.(states.StrValue))
				bytes := []byte(str)
				var iter func() (states.Value, bool, error)
				i := 0
				iter = func() (states.Value, bool, error) {
					if i >= len(bytes) {
						return nil, false, nil
					}
					v := states.NumValue(bytes[i])
					i++
					return v, true, nil
				}
				return states.ThunkFromIter(iter)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.StrType{},
			"split",
			nil,
			&types.ArrType{types.StrType{}},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				str := string(inputState.Value.(states.StrValue))
				fields := strings.Fields(str)
				var iter func() (states.Value, bool, error)
				i := 0
				iter = func() (states.Value, bool, error) {
					if i >= len(fields) {
						return nil, false, nil
					}
					v := states.StrValue(fields[i])
					i++
					return v, true, nil
				}
				return states.ThunkFromIter(iter)
			},
			nil,
		),
		expressions.RegularFuncer(
			types.StrType{},
			"split",
			[]*parameters.Parameter{
				parameters.SimpleParam(types.StrType{}),
			},
			&types.ArrType{types.StrType{}},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				str := string(inputState.Value.(states.StrValue))
				res0 := args[0](inputState, nil).Eval()
				if res0.Error != nil {
					return states.ThunkFromError(res0.Error)
				}
				sep := string(res0.Value.(states.StrValue))
				fields := strings.Split(str, sep)
				var iter func() (states.Value, bool, error)
				i := 0
				iter = func() (states.Value, bool, error) {
					if i >= len(fields) {
						return nil, false, nil
					}
					v := states.StrValue(fields[i])
					i++
					return v, true, nil
				}
				return states.ThunkFromIter(iter)
			},
			nil,
		),
		expressions.SimpleFuncer(
			&types.ArrType{types.StrType{}},
			"join",
			nil,
			types.StrType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				iter := states.IterFromValue(inputValue)
				buffer := bytes.Buffer{}
				for {
					value, ok, err := iter()
					if err != nil {
						return nil, err
					}
					if !ok {
						return states.StrValue(buffer.String()), nil
					}
					buffer.WriteString(string(value.(states.StrValue)))
				}
			},
		),
		expressions.SimpleFuncer(
			&types.ArrType{types.StrType{}},
			"join",
			[]types.Type{types.StrType{}},
			types.StrType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				iter := states.IterFromValue(inputValue)
				sep := string(argumentValues[0].(states.StrValue))
				buffer := bytes.Buffer{}
				firstWritten := false
				for {
					value, ok, err := iter()
					if err != nil {
						return nil, err
					}
					if !ok {
						return states.StrValue(buffer.String()), nil
					}
					if firstWritten {
						buffer.WriteString(sep)
					}
					buffer.WriteString(string(value.(states.StrValue)))
					firstWritten = true
				}
			},
		),
		expressions.SimpleFuncer(
			types.StrType{},
			"==",
			[]types.Type{types.StrType{}},
			types.BoolType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				str1 := string(inputValue.(states.StrValue))
				str2 := string(argumentValues[0].(states.StrValue))
				return states.BoolValue(str1 == str2), nil
			},
		),
		expressions.SimpleFuncer(
			types.StrType{},
			"+",
			[]types.Type{types.StrType{}},
			types.StrType{},
			func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
				str1 := string(inputValue.(states.StrValue))
				str2 := string(argumentValues[0].(states.StrValue))
				return states.StrValue(str1 + str2), nil
			},
		),
	})
}
