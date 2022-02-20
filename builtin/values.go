package builtin

import (
	"encoding/json"
	"strconv"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initValues() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.Funcer{
		expressions.SimpleFuncer(
			types.NewVar("A", types.Any{}),
			"id",
			nil,
			types.NewVar("A", types.Any{}),
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
		expressions.RegularFuncer(
			types.Str{},
			"parseJSON",
			nil,
			types.Any{},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				str := inputState.Value.(states.StrValue)
				var data interface{}
				err := json.Unmarshal([]byte(str), &data)
				if err != nil {
					return states.ThunkFromError(err)
				}
				return thunkFromData(data)
			},
			nil,
		),
	})
}

func thunkFromData(data interface{}) *states.Thunk {
	switch data := data.(type) {
	case nil:
		return states.ThunkFromValue(states.NullValue{})
	case bool:
		return states.ThunkFromValue(states.BoolValue(data))
	case float64:
		return states.ThunkFromValue(states.NumValue(data))
	case string:
		return states.ThunkFromValue(states.StrValue(data))
	case []interface{}:
		i := 0
		iter := func() (states.Value, bool, error) {
			if i >= len(data) {
				return nil, false, nil
			}
			res := thunkFromData(data[i]).Eval()
			if res.Error != nil {
				return nil, false, res.Error
			}
			i += 1
			return res.Value, true, nil
		}
		return states.ThunkFromIter(iter)
	case map[string]interface{}:
		obj := make(map[string]*states.Thunk)
		for k, v := range data {
			obj[k] = thunkFromData(v)
		}
		return states.ThunkFromValue(states.ObjValue(obj))
	default:
		return states.ThunkFromError(errors.ValueError(
			errors.Message("encountered unexpected object while converting from JSON"),
		))
	}
}
