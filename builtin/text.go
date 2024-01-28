package builtin

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var TextFuncers = []shapes.Funcer{
	shapes.SimpleFuncer(
		types.Str{},
		"<",
		[]*params.Param{
			params.SimpleParam("other", types.Str{}),
		},
		types.Bool{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.BoolValue(str1 < str2), nil
		},
	),
	shapes.SimpleFuncer(
		types.Str{},
		">",
		[]*params.Param{
			params.SimpleParam("other", types.Str{}),
		},
		types.Bool{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.BoolValue(str1 > str2), nil
		},
	),
	shapes.SimpleFuncer(
		types.Str{},
		"<=",
		[]*params.Param{
			params.SimpleParam("other", types.Str{}),
		},
		types.Bool{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.BoolValue(str1 <= str2), nil
		},
	),
	shapes.SimpleFuncer(
		types.Str{},
		">=",
		[]*params.Param{
			params.SimpleParam("other", types.Str{}),
		},
		types.Bool{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.BoolValue(str1 >= str2), nil
		},
	),
	shapes.SimpleFuncer(
		types.Str{},
		"+",
		[]*params.Param{
			params.SimpleParam("b", types.Str{}),
		},
		types.Str{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.StrValue(str1 + str2), nil
		},
	),

	shapes.Funcer{InputType: types.Str{}, Name: "bytes", Params: nil, OutputType: types.NewArr(types.Num{}), Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		str := inputState.Value.(states.StrValue)
		bytes := []byte(str)
		var output func() (states.Value, bool, error)
		i := 0
		output = func() (states.Value, bool, error) {
			if i >= len(bytes) {
				return nil, false, nil
			}
			v := states.NumValue(bytes[i])
			i++
			return v, true, nil
		}
		return states.ThunkFromIter(output)
	}, IDs: nil},

	shapes.Funcer{InputType: types.NewArr(types.Num{}), Name: "bytesToStr", Params: nil, OutputType: types.Str{}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		input := states.IterFromValue(inputState.Value)
		var output strings.Builder
		for {
			v, ok, err := input()
			if err != nil {
				return states.ThunkFromError(err)
			}
			if !ok {
				break
			}
			output.WriteByte(byte(v.(states.NumValue)))
		}
		return states.ThunkFromValue(states.StrValue(output.String()))
	}, IDs: nil},

	shapes.Funcer{InputType: types.Str{}, Name: "codePoints", Params: nil, OutputType: types.NewArr(types.Num{}), Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		str := inputState.Value.(states.StrValue)
		runes := []rune(str)
		i := 0
		output := func() (states.Value, bool, error) {
			if i >= len(runes) {
				return nil, false, nil
			}
			v := states.NumValue(runes[i])
			i++
			return v, true, nil
		}
		return states.ThunkFromIter(output)
	}, IDs: nil},

	shapes.Funcer{InputType: types.NewArr(types.Num{}), Name: "codePointsToStr", Params: nil, OutputType: types.Str{}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		input := states.IterFromValue(inputState.Value)
		var output strings.Builder
		for {
			v, ok, err := input()
			if err != nil {
				return states.ThunkFromError(err)
			}
			if !ok {
				break
			}
			output.WriteRune(rune(v.(states.NumValue)))
		}
		return states.ThunkFromValue(states.StrValue(output.String()))
	}, IDs: nil},

	shapes.Funcer{InputType: types.Str{}, Name: "fields", Params: nil, OutputType: types.NewArr(types.Str{}), Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		str := string(inputState.Value.(states.StrValue))
		fields := strings.Fields(str)
		i := 0
		iter := func() (states.Value, bool, error) {
			if i >= len(fields) {
				return nil, false, nil
			}
			v := states.StrValue(fields[i])
			i++
			return v, true, nil
		}
		return states.ThunkFromIter(iter)
	}, IDs: nil},

	shapes.Funcer{InputType: types.Str{}, Name: "indexOf", Params: []*params.Param{
		params.SimpleParam("needle", types.Str{}),
	}, OutputType: types.Num{}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		haystack := string(inputState.Value.(states.StrValue))
		needle, err := args[0](inputState.Clear(), nil).EvalStr()
		if err != nil {
			return states.ThunkFromError(err)
		}
		result := states.NumValue(strings.Index(haystack, needle))
		return states.ThunkFromValue(result)
	}, IDs: nil},

	// for Arr<Str> join Str
	shapes.SimpleFuncer(
		types.NewArr(types.Str{}),
		"join",
		nil,
		types.Str{},
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
	shapes.SimpleFuncer(
		types.NewArr(types.Str{}),
		"join",
		[]*params.Param{
			params.SimpleParam("glue", types.Str{}),
		},
		types.Str{},
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
	shapes.SimpleFuncer(
		types.Str{},
		"padEnd",
		[]*params.Param{
			params.SimpleParam("targetLength", types.Num{}),
			params.SimpleParam("padString", types.Str{}),
		},
		types.Str{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str := string(inputValue.(states.StrValue))
			length := int(argumentValues[0].(states.NumValue))
			padding := string(argumentValues[1].(states.StrValue))
			var builder strings.Builder
			builder.WriteString(str)
			for {
				delta := length - builder.Len()
				if delta <= 0 {
					break
				}
				if delta < len(padding) {
					builder.WriteString(padding[:delta])
					break
				}
				builder.WriteString(padding)
			}
			return states.StrValue(builder.String()), nil
		},
	),
	shapes.SimpleFuncer(
		types.Str{},
		"padStart",
		[]*params.Param{
			params.SimpleParam("targetLength", types.Num{}),
			params.SimpleParam("padString", types.Str{}),
		},
		types.Str{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str := string(inputValue.(states.StrValue))
			length := int(argumentValues[0].(states.NumValue))
			padding := string(argumentValues[1].(states.StrValue))
			var builder strings.Builder
			for {
				delta := length - len(str) - builder.Len()
				if delta <= 0 {
					break
				}
				if delta < len(padding) {
					builder.WriteString(padding[:delta])
					break
				}
				builder.WriteString(padding)
			}
			builder.WriteString(str)
			return states.StrValue(builder.String()), nil
		},
	),

	shapes.Funcer{InputType: types.Str{}, Name: "replaceFirst", Params: []*params.Param{
		params.SimpleParam("needle", types.Str{}),
		params.SimpleParam("replacement", types.Str{}),
	}, OutputType: types.Str{}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		haystack := inputState.Value.(states.StrValue)
		needle, err := args[0](inputState.Clear(), nil).EvalStr()
		if err != nil {
			return states.ThunkFromError(err)
		}
		replacement, err := args[1](inputState.Clear(), nil).EvalStr()
		if err != nil {
			return states.ThunkFromError(err)
		}
		result := strings.Replace(string(haystack), string(needle), string(replacement), 1)
		return states.ThunkFromValue(states.StrValue(result))
	}, IDs: nil},

	shapes.Funcer{InputType: types.Str{}, Name: "replaceAll", Params: []*params.Param{
		params.SimpleParam("needle", types.Str{}),
		params.SimpleParam("replacement", types.Str{}),
	}, OutputType: types.Str{}, Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
		haystack := inputState.Value.(states.StrValue)
		needle, err := args[0](inputState.Clear(), nil).EvalStr()
		if err != nil {
			return states.ThunkFromError(err)
		}
		replacement, err := args[1](inputState.Clear(), nil).EvalStr()
		if err != nil {
			return states.ThunkFromError(err)
		}
		result := strings.ReplaceAll(string(haystack), string(needle), string(replacement))
		return states.ThunkFromValue(states.StrValue(result))
	}, IDs: nil},

	shapes.SimpleFuncer(
		types.Str{},
		"startsWith",
		[]*params.Param{
			params.SimpleParam("needle", types.Str{}),
		},
		types.Bool{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.BoolValue(strings.HasPrefix(str1, str2)), nil
		},
	),
	shapes.SimpleFuncer(
		types.Str{},
		"endsWith",
		[]*params.Param{
			params.SimpleParam("needle", types.Str{}),
		},
		types.Bool{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.BoolValue(strings.HasSuffix(str1, str2)), nil
		},
	),
	shapes.SimpleFuncer(
		types.Str{},
		"slice",
		[]*params.Param{
			params.SimpleParam("start", types.Num{}),
		},
		types.Str{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str := string(inputValue.(states.StrValue))
			start := int(argumentValues[0].(states.NumValue))
			if start < 0 {
				start = len(str) + start
				if start < 0 {
					start = 0
				}
			}
			return states.StrValue(str[start:]), nil
		},
	),
	shapes.SimpleFuncer(
		types.Str{},
		"slice",
		[]*params.Param{
			params.SimpleParam("start", types.Num{}),
			params.SimpleParam("end", types.Num{}),
		},
		types.Str{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str := string(inputValue.(states.StrValue))
			start := int(argumentValues[0].(states.NumValue))
			end := int(argumentValues[1].(states.NumValue))
			if start < 0 {
				start = len(str) + start
			}
			if start < 0 {
				start = 0
			}
			if end < 0 {
				end = len(str) + end
			}
			if end < start {
				end = start
			}
			return states.StrValue(str[start:end]), nil
		},
	),
	shapes.SimpleFuncer(
		types.Str{},
		"repeat",
		[]*params.Param{
			params.SimpleParam("times", types.Num{}),
		},
		types.Str{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str := string(inputValue.(states.StrValue))
			n := int(argumentValues[0].(states.NumValue))
			if n < 0 {
				n = 0
			}
			return states.StrValue(strings.Repeat(str, n)), nil
		},
	),
	// for Str trim Str
	shapes.SimpleFuncer(
		types.Str{},
		"trim",
		nil,
		types.Str{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str := string(inputValue.(states.StrValue))
			return states.StrValue(strings.TrimSpace(str)), nil
		},
	),
	// for Str trimStart Str
	shapes.SimpleFuncer(
		types.Str{},
		"trimStart",
		nil,
		types.Str{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str := string(inputValue.(states.StrValue))
			return states.StrValue(strings.TrimLeftFunc(str, unicode.IsSpace)), nil
		},
	),
	// for Str trimEnd Str
	shapes.SimpleFuncer(
		types.Str{},
		"trimEnd",
		nil,
		types.Str{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str := string(inputValue.(states.StrValue))
			return states.StrValue(strings.TrimRightFunc(str, unicode.IsSpace)), nil
		},
	),
}
