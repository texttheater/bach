package builtin

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/expressions"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var TextFuncers = []expressions.Funcer{
	// for Str <(Str) Bool
	expressions.SimpleFuncer(
		types.Str{},
		"<",
		[]types.Type{types.Str{}},
		types.Bool{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.BoolValue(str1 < str2), nil
		},
	),
	// for Str >(Str) Bool
	expressions.SimpleFuncer(
		types.Str{},
		">",
		[]types.Type{types.Str{}},
		types.Bool{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.BoolValue(str1 > str2), nil
		},
	),
	// for Str <=(Str) Bool
	expressions.SimpleFuncer(
		types.Str{},
		"<=",
		[]types.Type{types.Str{}},
		types.Bool{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.BoolValue(str1 <= str2), nil
		},
	),
	// for Str >=(Str) Bool
	expressions.SimpleFuncer(
		types.Str{},
		">=",
		[]types.Type{types.Str{}},
		types.Bool{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.BoolValue(str1 >= str2), nil
		},
	),
	// for Str +(Str) Bool
	expressions.SimpleFuncer(
		types.Str{},
		"+",
		[]types.Type{types.Str{}},
		types.Str{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.StrValue(str1 + str2), nil
		},
	),
	// for Str bytes Arr<Num>
	expressions.RegularFuncer(
		types.Str{},
		"bytes",
		nil,
		types.NewArr(types.Num{}),
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		},
		nil,
	),
	// for Arr<Num> bytesToStr Str
	expressions.RegularFuncer(
		types.NewArr(types.Num{}),
		"bytesToStr",
		nil,
		types.Str{},
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		},
		nil,
	),
	// for Str codePoints Arr<Num>
	expressions.RegularFuncer(
		types.Str{},
		"codePoints",
		nil,
		types.NewArr(types.Num{}),
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		},
		nil,
	),
	// for Arr<Num> codePointsToStr Str
	expressions.RegularFuncer(
		types.NewArr(types.Num{}),
		"codePointsToStr",
		nil,
		types.Str{},
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		},
		nil,
	),
	// for Str fields Arr<Str>
	expressions.RegularFuncer(
		types.Str{},
		"fields",
		nil,
		types.NewArr(types.Str{}),
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		},
		nil,
	),
	// for Str indexOf(Str) Num
	expressions.RegularFuncer(
		types.Str{},
		"indexOf",
		[]*params.Param{
			params.SimpleParam("needle", types.Str{}),
		},
		types.Num{},
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			haystack := string(inputState.Value.(states.StrValue))
			needle, err := args[0](inputState.Clear(), nil).EvalStr()
			if err != nil {
				return states.ThunkFromError(err)
			}
			result := states.NumValue(strings.Index(haystack, needle))
			return states.ThunkFromValue(result)
		},
		nil,
	),
	// for Arr<Str> join Str
	expressions.SimpleFuncer(
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
	// for Arr<Str> join(Str) Str
	expressions.SimpleFuncer(
		types.NewArr(types.Str{}),
		"join",
		[]types.Type{types.Str{}},
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
	// for Str padEnd(Num, Str) Str
	expressions.SimpleFuncer(
		types.Str{},
		"padEnd",
		[]types.Type{types.Num{}, types.Str{}},
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
	// for Str padEnd(Num, Str) Str
	expressions.SimpleFuncer(
		types.Str{},
		"padStart",
		[]types.Type{types.Num{}, types.Str{}},
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
	// for Str replaceFirst(Str, Str) Str
	expressions.RegularFuncer(
		types.Str{},
		"replaceFirst",
		[]*params.Param{
			params.SimpleParam("needle", types.Str{}),
			params.SimpleParam("replacement", types.Str{}),
		},
		types.Str{},
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		},
		nil,
	),
	// for Str replaceAll(Str, Str) Str
	expressions.RegularFuncer(
		types.Str{},
		"replaceAll",
		[]*params.Param{
			params.SimpleParam("needle", types.Str{}),
			params.SimpleParam("replacement", types.Str{}),
		},
		types.Str{},
		func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		},
		nil,
	),
	// for Str startsWith(Str) Bool
	expressions.SimpleFuncer(
		types.Str{},
		"startsWith",
		[]types.Type{types.Str{}},
		types.Bool{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.BoolValue(strings.HasPrefix(str1, str2)), nil
		},
	),
	// for Str endsWith(Str) Bool
	expressions.SimpleFuncer(
		types.Str{},
		"endsWith",
		[]types.Type{types.Str{}},
		types.Bool{},
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.BoolValue(strings.HasSuffix(str1, str2)), nil
		},
	),
	// for Str slice(Num) Str
	expressions.SimpleFuncer(
		types.Str{},
		"slice",
		[]types.Type{types.Num{}},
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
	// for Str slice(Num, Num) Str
	expressions.SimpleFuncer(
		types.Str{},
		"slice",
		[]types.Type{
			types.Num{},
			types.Num{},
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
	// for Str repeat(Num) Str
	expressions.SimpleFuncer(
		types.Str{},
		"repeat",
		[]types.Type{
			types.Num{},
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
	expressions.SimpleFuncer(
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
	expressions.SimpleFuncer(
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
	expressions.SimpleFuncer(
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
