package builtin

import (
	"encoding/json"
	"strconv"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var ValueFuncers = []shapes.Funcer{

	shapes.Funcer{InputType: types.Any{}, Name: "==", Params: []*params.Param{
		params.SimpleParam("other", (types.Any{})),
	}, OutputType: types.Bool{}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		a := inputState.Value
		b, err := args[0](inputState.Clear(), nil).Eval()
		if err != nil {
			return states.ThunkFromError(err)
		}
		equal, err := a.Equal(b)
		if err != nil {
			return states.ThunkFromError(err)
		}
		return states.ThunkFromValue(states.BoolValue(equal))
	}, IDs: nil},

	// for <A> id <A>
	shapes.SimpleFuncer(
		types.NewVar("A", types.Any{}),
		"id",
		nil,
		types.NewVar("A", types.Any{}),
		func(inputValue states.Value, argValues []states.Value) (states.Value, error) {
			return inputValue, nil
		},
	),

	shapes.Funcer{InputType: types.Str{}, Name: "parseFloat", Params: nil, OutputType: types.Num{}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		s := string(inputState.Value.(states.StrValue))
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return states.ThunkFromError(errors.ValueError(
				errors.Pos(pos),
				errors.Code(errors.UnexpectedValue),
				errors.Message(err.Error()),
			))
		}
		return states.ThunkFromValue(states.NumValue(n))
	}, IDs: nil},

	shapes.Funcer{InputType: types.Str{}, Name: "parseInt", Params: []*params.Param{
		params.SimpleParam("base", types.Num{}),
	}, OutputType: types.Num{}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		s := string(inputState.Value.(states.StrValue))
		b, err := args[0](inputState.Clear(), nil).EvalNum()
		if err != nil {
			return states.ThunkFromError(err)
		}
		n, err := strconv.ParseInt(s, int(b), 64)
		if err != nil {
			return states.ThunkFromError(errors.ValueError(
				errors.Pos(pos),
				errors.Code(errors.UnexpectedValue),
				errors.Message(err.Error()),
			))
		}
		return states.ThunkFromValue(states.NumValue(n))
	}, IDs: nil},

	shapes.Funcer{InputType: types.Str{}, Name: "parseInt", Params: nil, OutputType: types.Num{}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		s := string(inputState.Value.(states.StrValue))
		b := 10
		n, err := strconv.ParseInt(s, int(b), 64)
		if err != nil {
			return states.ThunkFromError(errors.ValueError(
				errors.Pos(pos),
				errors.Code(errors.UnexpectedValue),
				errors.Message(err.Error()),
			))
		}
		return states.ThunkFromValue(states.NumValue(n))
	}, IDs: nil},

	shapes.Funcer{InputType: types.Str{}, Name: "parseJSON", Params: nil, OutputType: types.Any{}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		str := inputState.Value.(states.StrValue)
		var data any
		err := json.Unmarshal([]byte(str), &data)
		if err != nil {
			return states.ThunkFromError(errors.ValueError(
				errors.Pos(pos),
				errors.Code(errors.UnexpectedValue),
				errors.Message(err.Error()),
			))
		}
		return thunkFromData(data, pos)
	}, IDs: nil},

	shapes.Funcer{InputType: types.Any{}, Name: "toJSON", Params: nil, OutputType: types.Str{}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		data, err := inputState.Value.Data()
		if err != nil {
			return states.ThunkFromError(err)
		}
		bytes, err := json.Marshal(data)
		if err != nil {
			return states.ThunkFromError(errors.ValueError(
				errors.Pos(pos),
				errors.Code(errors.UnexpectedValue),
				errors.Message(err.Error()),
			))
		}
		return states.ThunkFromValue(states.StrValue(bytes))
	}, IDs: nil},

	shapes.Funcer{InputType: types.NewArr(types.NewTup([]types.Type{
		types.Str{},
		types.NewVar("A", types.Any{}),
	})), Name: "toObj", Params: nil, OutputType: types.Obj{
		Props: map[string]types.Type{},
		Rest:  types.NewVar("A", types.Any{}),
	}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		var res states.ObjValue = make(map[string]*states.Thunk)
		iter := states.IterFromValue(inputState.Value)
		for {
			val, ok, err := iter()
			if err != nil {
				return states.ThunkFromError(err)
			}
			if !ok {
				break
			}
			arr := val.(*states.ArrValue)
			prop := string(arr.Head.(states.StrValue))
			tail, err := arr.Tail.Eval()
			if err != nil {
				return states.ThunkFromError(err)
			}
			v := states.ThunkFromValue(tail.(*states.ArrValue).Head)
			res[prop] = v
		}
		return states.ThunkFromValue(res)
	}, IDs: nil},

	shapes.Funcer{InputType: types.Any{}, Name: "toStr", Params: nil, OutputType: types.Str{}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		str, err := inputState.Value.Str()
		if err != nil {
			return states.ThunkFromError(err)
		}
		return states.ThunkFromValue(states.StrValue(str))
	}, IDs: nil},
}

func thunkFromData(data any, pos lexer.Position) *states.Thunk {
	switch data := data.(type) {
	case nil:
		return states.ThunkFromValue(states.NullValue{})
	case bool:
		return states.ThunkFromValue(states.BoolValue(data))
	case float64:
		return states.ThunkFromValue(states.NumValue(data))
	case string:
		return states.ThunkFromValue(states.StrValue(data))
	case []any:
		i := 0
		iter := func() (states.Value, bool, error) {
			if i >= len(data) {
				return nil, false, nil
			}
			val, err := thunkFromData(data[i], pos).Eval()
			if err != nil {
				return nil, false, err
			}
			i += 1
			return val, true, nil
		}
		return states.ThunkFromIter(iter)
	case map[string]any:
		obj := make(map[string]*states.Thunk)
		for k, v := range data {
			obj[k] = thunkFromData(v, pos)
		}
		return states.ThunkFromValue(states.ObjValue(obj))
	default:
		return states.ThunkFromError(errors.ValueError(
			errors.Pos(pos),
			errors.Code(errors.UnexpectedValue),
			errors.Message("encountered unexpected object while converting from JSON"),
		))
	}
}
