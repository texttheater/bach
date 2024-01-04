package builtin

import (
	"encoding/json"
	"strconv"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func initValues() {
	InitialShape.Stack = InitialShape.Stack.PushAll([]expressions.FuncerDefinition{
		// for Any ==(Any) Bool
		expressions.RegularFuncer(
			types.Any{},
			"==",
			[]*params.Param{
				params.SimpleParam((types.Any{})),
			},
			types.Bool{},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
			},
			nil,
		),
		// for <A> id <A>
		expressions.SimpleFuncer(
			types.NewVar("A", types.Any{}),
			"id",
			nil,
			types.NewVar("A", types.Any{}),
			func(inputValue states.Value, argValues []states.Value) (states.Value, error) {
				return inputValue, nil
			},
		),
		// for Str parseFloat Num
		expressions.RegularFuncer(
			types.Str{},
			"parseFloat",
			nil,
			types.Num{},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
			},
			nil,
		),
		// for Str parseInt(Num) Num
		expressions.RegularFuncer(
			types.Str{},
			"parseInt",
			[]*params.Param{
				params.SimpleParam(types.Num{}),
			},
			types.Num{},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
			},
			nil,
		),
		// for Str parseInt Num
		expressions.RegularFuncer(
			types.Str{},
			"parseInt",
			nil,
			types.Num{},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
			},
			nil,
		),
		// for Str parseJSON Any
		expressions.RegularFuncer(
			types.Str{},
			"parseJSON",
			nil,
			types.Any{},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
			},
			nil,
		),
		// for Any toJSON Str
		expressions.RegularFuncer(
			types.Any{},
			"toJSON",
			nil,
			types.Str{},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
			},
			nil,
		),
		// for Arr<Tup<Str, <A>>> toObj Obj<<A>>
		expressions.RegularFuncer(
			types.NewArr(types.NewTup([]types.Type{
				types.Str{},
				types.NewVar("A", types.Any{}),
			})),
			"toObj",
			nil,
			types.Obj{
				Props: map[string]types.Type{},
				Rest:  types.NewVar("A", types.Any{}),
			},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
			},
			nil,
		),
		// for Any toStr Str
		expressions.RegularFuncer(
			types.Any{},
			"toStr",
			nil,
			types.Str{},
			func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
				str, err := inputState.Value.Str()
				if err != nil {
					return states.ThunkFromError(err)
				}
				return states.ThunkFromValue(states.StrValue(str))
			},
			nil,
		),
	})
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
