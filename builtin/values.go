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
	shapes.Funcer{
		Summary:          "Checks two values for equality.",
		InputType:        types.AnyType{},
		InputDescription: "a value",
		Name:             "==",
		Params: []*params.Param{
			params.SimpleParam("other", "another value", (types.AnyType{})),
		},
		OutputType:        types.BoolType{},
		OutputDescription: "true if the input is the same as other, false otherwise",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`null ==null`, `Bool`, `true`, nil},
			{`null =={}`, `Bool`, `false`, nil},
			{`true ==true`, `Bool`, `true`, nil},
			{`true ==false`, `Bool`, `false`, nil},
			{`true ==[]`, `Bool`, `false`, nil},
			{`1 ==1.0`, `Bool`, `true`, nil},
			{`1 ==2`, `Bool`, `false`, nil},
			{`57 =="a"`, `Bool`, `false`, nil},
			{`"abc" =="abc"`, `Bool`, `true`, nil},
			{`"" =="abc"`, `Bool`, `false`, nil},
			{`"" ==null`, `Bool`, `false`, nil},
			{`[false, 1.0, "ab"] ==[false, 1, "a" +"b"]`, `Bool`, `true`, nil},
			{`[] ==[11]`, `Bool`, `false`, nil},
			{`["a"] =={a: 1}`, `Bool`, `false`, nil},
			{`{a: 1, b: 2} =={b: 2, a: 1}`, `Bool`, `true`, nil},
			{`{a: 1, b: 2} =={a: 2, b: 1}`, `Bool`, `false`, nil},
			{`{} ==[]`, `Bool`, `false`, nil},
		},
	},
	shapes.SimpleFuncer(
		"The identity function.",
		types.NewTypeVar("A", types.AnyType{}),
		"any value",
		"id",
		nil,
		types.NewTypeVar("A", types.AnyType{}),
		"the input value",
		"",
		func(inputValue states.Value, argValues []states.Value) (states.Value, error) {
			return inputValue, nil
		},
		[]shapes.Example{
			{`123 id`, `Num`, `123`, nil},
			{`"abc" id`, `Str`, `"abc"`, nil},
			{`false if id then 1 else fatal ok`, `Num`, `"abc"`, errors.ValueError(
				errors.Code(errors.UnexpectedValue),
				errors.GotValue(states.BoolValue(false)),
			)},
		},
	),
	shapes.Funcer{InputType: types.StrType{},
		Summary:           "Parses the string representation of a floating-point number.",
		Name:              "parseFloat",
		InputDescription:  "a floating-point number in string representation",
		Params:            nil,
		OutputType:        types.NumType{},
		OutputDescription: "the corresponding Num value",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			s := string(inputState.Value.(states.StrValue))
			n, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return states.ThunkFromError(errors.ValueError(
					errors.Pos(pos),
					errors.Code(errors.UnexpectedValue),
					errors.GotValue(inputState.Value.(states.StrValue)),
					errors.Message(err.Error()),
				))
			}
			return states.ThunkFromValue(states.NumValue(n))
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`"4.567" parseFloat`, `Num`, `4.567`, nil},
			{`"4.567e3" parseFloat`, `Num`, `4567`, nil},
			{`"4.567abcdefgh" parseFloat`, `Num`, ``, errors.ValueError(
				errors.Code(errors.UnexpectedValue),
				errors.GotValue(states.StrValue("4.567abcdefgh")),
			)},
		},
	},

	shapes.Funcer{
		Summary:          "Parses the string representation of an integer number.",
		InputType:        types.StrType{},
		InputDescription: "an integer number in string representation",
		Name:             "parseInt",
		Params: []*params.Param{
			params.SimpleParam("base", "the base that the input is in", types.NumType{}),
		},
		OutputType:        types.NumType{},
		OutputDescription: "the corresponding Num value",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
					errors.GotValue(inputState.Value.(states.StrValue)),
					errors.Message(err.Error()),
				))
			}
			return states.ThunkFromValue(states.NumValue(n))
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`"123" parseInt(10)`, `Num`, `123`, nil},
			{`"ff" parseInt(10)`, `Num`, ``, errors.ValueError(
				errors.Code(errors.UnexpectedValue),
				errors.GotValue(states.StrValue("ff")),
			)},
			{`"ff" parseInt(16)`, `Num`, `255`, nil},
			{`"0xFF" parseInt(16)`, `Num`, `255`, errors.ValueError(
				errors.Code(errors.UnexpectedValue),
				errors.GotValue(states.StrValue("0xFF")),
			)},
		},
	},
	shapes.Funcer{
		Summary:           "Parses the string representation of a base-10 integer number.",
		InputType:         types.StrType{},
		InputDescription:  "an integer number in base-10 string representation",
		Name:              "parseInt",
		Params:            nil,
		OutputType:        types.NumType{},
		OutputDescription: "the corresponding Num value",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			s := string(inputState.Value.(states.StrValue))
			b := 10
			n, err := strconv.ParseInt(s, int(b), 64)
			if err != nil {
				return states.ThunkFromError(errors.ValueError(
					errors.Pos(pos),
					errors.Code(errors.UnexpectedValue),
					errors.GotValue(inputState.Value.(states.StrValue)),
					errors.Message(err.Error()),
				))
			}
			return states.ThunkFromValue(states.NumValue(n))
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`"123" parseInt`, `Num`, `123`, nil},
			{`"   123 " parseInt`, `Num`, `123`, errors.ValueError(
				errors.Code(errors.UnexpectedValue),
				errors.GotValue(states.StrValue("   123 ")),
			)},
			{`"077" parseInt`, `Num`, `77`, nil},
			{`"1.9" parseInt`, `Num`, `1`, errors.ValueError(
				errors.Code(errors.UnexpectedValue),
				errors.GotValue(states.StrValue("1.9")),
			)},
			{`"ff" parseInt`, `Num`, `16`, errors.ValueError(
				errors.Code(errors.UnexpectedValue),
				errors.GotValue(states.StrValue("ff")),
			)},
			{`"xyz" parseInt`, `Num`, `16`, errors.ValueError(
				errors.Code(errors.UnexpectedValue),
				errors.GotValue(states.StrValue("xyz")),
			)},
		},
	},

	shapes.Funcer{InputType: types.StrType{}, Name: "parseJSON", Params: nil, OutputType: types.AnyType{}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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

	shapes.Funcer{InputType: types.AnyType{}, Name: "toJSON", Params: nil, OutputType: types.StrType{}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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

	shapes.Funcer{InputType: types.NewArrType(types.NewTup([]types.Type{
		types.StrType{},
		types.NewTypeVar("A", types.AnyType{}),
	})), Name: "toObj", Params: nil, OutputType: types.ObjType{
		Props: map[string]types.Type{},
		Rest:  types.NewTypeVar("A", types.AnyType{}),
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

	shapes.Funcer{InputType: types.AnyType{}, Name: "toStr", Params: nil, OutputType: types.StrType{}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
