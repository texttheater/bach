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
		"Compares strings lexicographically.",
		types.StrType{},
		"a string",
		"<",
		[]*params.Param{
			params.SimpleParam("other", "another string", types.StrType{}),
		},
		types.BoolType{},
		"true if the input appears before other in lexicographical order, false otherwise",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.BoolValue(str1 < str2), nil
		},
		[]shapes.Example{
			{`"a" <"b"`, `Bool`, `true`, nil},
			{`"Ab" <"A"`, `Bool`, `false`, nil},
			{`"cd" <"cd"`, `Bool`, `false`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Compares strings lexicographically.",
		types.StrType{},
		"a string",
		">",
		[]*params.Param{
			params.SimpleParam("other", "another string", types.StrType{}),
		},
		types.BoolType{},
		"true if the input appears after other in lexicographical order, false otherwise",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.BoolValue(str1 > str2), nil
		},
		[]shapes.Example{
			{`"a" >"b"`, `Bool`, `false`, nil},
			{`"Ab" >"A"`, `Bool`, `true`, nil},
			{`"cd" >"cd"`, `Bool`, `false`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Compares strings lexicographically.",
		types.StrType{},
		"a string",
		"<=",
		[]*params.Param{
			params.SimpleParam("other", "another string", types.StrType{}),
		},
		types.BoolType{},
		"true if the input appears before other in lexicographical order or is equal to it, false otherwise",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.BoolValue(str1 <= str2), nil
		},
		[]shapes.Example{
			{`"a" <="b"`, `Bool`, `true`, nil},
			{`"Ab" <="A"`, `Bool`, `false`, nil},
			{`"cd" <="cd"`, `Bool`, `true`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Compares strings lexicographically.",
		types.StrType{},
		"a string",
		">=",
		[]*params.Param{
			params.SimpleParam("other", "another string", types.StrType{}),
		},
		types.BoolType{},
		"true if the input appears after other in lexicographical order or is equal to it, false otherwise",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.BoolValue(str1 >= str2), nil
		},
		[]shapes.Example{
			{`"a" >="b"`, `Bool`, `false`, nil},
			{`"Ab" >="A"`, `Bool`, `true`, nil},
			{`"cd" >="cd"`, `Bool`, `true`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Concatenates two strings.",
		types.StrType{},
		"a string",
		"+",
		[]*params.Param{
			params.SimpleParam("b", "another string", types.StrType{}),
		},
		types.StrType{},
		"The input and b, concatenated.",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.StrValue(str1 + str2), nil
		},
		[]shapes.Example{
			{`"ab" +"cd"`, `Str`, `"abcd"`, nil},
			{`"ab" +""`, `Str`, `"ab"`, nil},
			{`"" +"cd"`, `Str`, `"cd"`, nil},
		},
	),
	shapes.Funcer{
		Summary:           "Converts a string to bytes.",
		InputType:         types.StrType{},
		InputDescription:  "a string",
		Name:              "bytes",
		Params:            nil,
		OutputType:        types.NewArrType(types.NumType{}),
		OutputDescription: "The UTF-8 bytes representing the string.",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`"abc" bytes`, `Arr<Num...>`, `[97, 98, 99]`, nil},
			{`"Köln" bytes`, `Arr<Num...>`, `[75, 195, 182, 108, 110]`, nil},
			{`"日本語" bytes`, `Arr<Num...>`, `[230, 151, 165, 230, 156, 172, 232, 170, 158]`, nil},
			{`"\x00" bytes`, `Arr<Num...>`, `[0]`, nil},
		},
	},
	shapes.Funcer{
		Summary:           "Converts bytes to a string.",
		InputType:         types.NewArrType(types.NumType{}),
		InputDescription:  "an array of numbers (interpreted modulo 256 as UTF-8 bytes)",
		Name:              "bytesToStr",
		Params:            nil,
		OutputType:        types.StrType{},
		OutputDescription: "the string represented by the input",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`[97, 98, 99] bytesToStr`, `Str`, `"abc"`, nil},
			{`[75, 195, 182, 108, 110] bytesToStr`, `Str`, `"Köln"`, nil},
			{`[230, 151, 165, 230, 156, 172, 232, 170, 158] bytesToStr`, `Str`, `"日本語"`, nil},
			{`[0] bytesToStr`, `Str`, `"\x00"`, nil},
			{`[256] bytesToStr`, `Str`, `"\x00"`, nil},
		},
	},
	shapes.Funcer{
		Summary:           "Converts a string to Unicode code points.",
		InputType:         types.StrType{},
		InputDescription:  "a string",
		Name:              "codePoints",
		Params:            nil,
		OutputType:        types.NewArrType(types.NumType{}),
		OutputDescription: "the input represented as a sequence of code points",
		Notes:             "If the input string contains invalid UTF-8 byte sequences, they will be represented by the Unicode replacement character (code point 65533).",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`"abc" codePoints`, `Arr<Num...>`, `[97, 98, 99]`, nil},
			{`"Köln" codePoints`, `Arr<Num...>`, `[75, 246, 108, 110]`, nil},
			{`"日本語" codePoints`, `Arr<Num...>`, `[26085, 26412, 35486]`, nil},
			{`"\x80" codePoints`, `Arr<Num...>`, `[65533]`, nil},
		},
	},
	shapes.Funcer{
		Summary:           "Converts Unicode code points to a string.",
		InputType:         types.NewArrType(types.NumType{}),
		InputDescription:  "a sequence of numbers",
		Name:              "codePointsToStr",
		Params:            nil,
		OutputType:        types.StrType{},
		OutputDescription: "UTF-8 encoded version of the input",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`[97, 98, 99] codePointsToStr`, `Str`, `"abc"`, nil},
			{`[75, 246, 108, 110] codePointsToStr`, `Str`, `"Köln"`, nil},
			{`[26085, 26412, 35486] codePointsToStr`, `Str`, `"日本語"`, nil},
			{`[65533] codePointsToStr`, `Str`, `"�"`, nil},
		},
	},
	shapes.Funcer{
		Summary:           "Splits a string around whitespace.",
		InputType:         types.StrType{},
		InputDescription:  "a string",
		Name:              "fields",
		Params:            nil,
		OutputType:        types.NewArrType(types.StrType{}),
		OutputDescription: "the result of splitting the string around any kind or amount of white space",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`"  foo bar  baz   " fields`, `Arr<Str...>`, `["foo", "bar", "baz"]`, nil},
		},
	},
	shapes.Funcer{
		Summary:          "Finds the position of a string within another.",
		InputType:        types.StrType{},
		InputDescription: "string to search inside",
		Name:             "indexOf",
		Params: []*params.Param{
			params.SimpleParam("needle", "string to search for", types.StrType{}),
		},
		OutputType:        types.NumType{},
		OutputDescription: "offset of first occurrence of needle from the beginning of the input, measured in bytes, or -1 if none",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			haystack := string(inputState.Value.(states.StrValue))
			needle, err := args[0](inputState.Clear(), nil).EvalStr()
			if err != nil {
				return states.ThunkFromError(err)
			}
			result := states.NumValue(strings.Index(haystack, needle))
			return states.ThunkFromValue(result)
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`"abc" indexOf("bc")`, `Num`, `1`, nil},
			{`"abc" indexOf("a")`, `Num`, `0`, nil},
			{`"abc" indexOf("d")`, `Num`, `-1`, nil},
			{`"Köln" indexOf("l")`, `Num`, `3`, nil},
		},
	},
	shapes.SimpleFuncer(
		"Concatenates any number of strings.",
		types.NewArrType(types.StrType{}),
		"an array of strings",
		"join",
		nil,
		types.StrType{},
		"the concatenation of all the strings in the input",
		"",
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
		[]shapes.Example{
			{`["ab", "cd", "ef"] join`, `Str`, `"abcdef"`, nil},
			{`["ab", "cd"] join`, `Str`, `"abcd"`, nil},
			{`["ab"] join`, `Str`, `"ab"`, nil},
			{`for Any def f Arr<Str...> as [] ok f join`, `Str`, `""`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Concatenates any number of strings with a custom delimiter.",
		types.NewArrType(types.StrType{}),
		"an array of strings",
		"join",
		[]*params.Param{
			params.SimpleParam("glue", "delimiter to put between strings", types.StrType{}),
		},
		types.StrType{},
		"the concatenation of all the strings in the input, with the glue in between",
		"",
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
		[]shapes.Example{
			{`["ab", "cd", "ef"] join(";")`, `Str`, `"ab;cd;ef"`, nil},
			{`["ab", "", "cd"] join(";")`, `Str`, `"ab;;cd"`, nil},
			{`["ab"] join(";")`, `Str`, `"ab"`, nil},
			{`for Any def f Arr<Str...> as [] ok f join(";")`, `Str`, `""`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Right-pads a string.",
		types.StrType{},
		"a string",
		"padEnd",
		[]*params.Param{
			params.SimpleParam("targetLength", "minimum string length for the output", types.NumType{}),
			params.SimpleParam("padString", "string to use as padding to bring the input to the desired length", types.StrType{}),
		},
		types.StrType{},
		"the input followed by as many repetitions of padString as necessary to reach targetLength",
		"padString is usually one character but can be longer, in which case each repetition starts from the left.",
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
		[]shapes.Example{
			{`"z" padEnd(2, " ")`, `Str`, `"z "`, nil},
			{`"z" padEnd(3, " ")`, `Str`, `"z  "`, nil},
			{`"zzz" padEnd(3, " ")`, `Str`, `"zzz"`, nil},
			{`"zzzz" padEnd(3, " ")`, `Str`, `"zzzz"`, nil},
			{`"zzzz" padEnd(7, "ab")`, `Str`, `"zzzzaba"`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Left-pads a string.",
		types.StrType{},
		"a string",
		"padStart",
		[]*params.Param{
			params.SimpleParam("targetLength", "minimum string length for the output", types.NumType{}),
			params.SimpleParam("padString", "string to use as padding to bring the input to the desired length", types.StrType{}),
		},
		types.StrType{},
		"the input following as many repetitions of padString as necessary to reach targetLength",
		"padString is usually one character but can be longer, in which case each repetition starts from the left.",
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
		[]shapes.Example{
			{`"z" padStart(2, " ")`, `Str`, `" z"`, nil},
			{`"z" padStart(3, " ")`, `Str`, `"  z"`, nil},
			{`"zzz" padStart(3, " ")`, `Str`, `"zzz"`, nil},
			{`"zzzz" padStart(3, " ")`, `Str`, `"zzzz"`, nil},
			{`"zzzz" padStart(7, "ab")`, `Str`, `"abazzzz"`, nil},
		},
	),
	shapes.Funcer{
		Summary:          "Replaces substrings.",
		InputType:        types.StrType{},
		InputDescription: "a string",
		Name:             "replaceAll",
		Params: []*params.Param{
			params.SimpleParam("needle", "a substring to look for", types.StrType{}),
			params.SimpleParam("replacement", "a new string to replace needle with", types.StrType{}),
		},
		OutputType:        types.StrType{},
		OutputDescription: "the input with all occurrences of needle replaced with replacement",
		Notes:             "More precisely, whenever there are two or more overlapping occurrences of needle in the input, only the first one is replaced.",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`"ababa" replaceAll("b", "c")`, `Str`, `"acaca"`, nil},
		},
	},
	shapes.Funcer{
		Summary:          "Replaces a substring.",
		InputType:        types.StrType{},
		InputDescription: "a string",
		Name:             "replaceFirst",
		Params: []*params.Param{
			params.SimpleParam("needle", "a substring to look for", types.StrType{}),
			params.SimpleParam("replacement", "a new string to replace needle with", types.StrType{}),
		},
		OutputType:        types.StrType{},
		OutputDescription: "the input with the first ocurrence of needle replaced with replacement",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
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
		IDs: nil,
		Examples: []shapes.Example{
			{`"ababa" replaceFirst("b", "c")`, `Str`, `"acaba"`, nil},
		},
	},
	shapes.SimpleFuncer(
		"Checks whether a string starts with a specific substring.",
		types.StrType{},
		"a string",
		"startsWith",
		[]*params.Param{
			params.SimpleParam("needle", "a prefix to look for", types.StrType{}),
		},
		types.BoolType{},
		"true if the input starts with needle, false otherwise",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.BoolValue(strings.HasPrefix(str1, str2)), nil
		},
		[]shapes.Example{
			{`"abc" startsWith("ab")`, `Bool`, `true`, nil},
			{`"abc" startsWith("b")`, `Bool`, `false`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Checks whether a string ends with a specific substring.",
		types.StrType{},
		"a string",
		"endsWith",
		[]*params.Param{
			params.SimpleParam("needle", "a suffix to look for", types.StrType{}),
		},
		types.BoolType{},
		"true if the input ends with needle, false otherwise",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str1 := string(inputValue.(states.StrValue))
			str2 := string(argumentValues[0].(states.StrValue))
			return states.BoolValue(strings.HasSuffix(str1, str2)), nil
		},
		[]shapes.Example{
			{`"abc" endsWith("bc")`, `Bool`, `true`, nil},
			{`"abc" endsWith("b")`, `Bool`, `false`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Concatenates multiple repetitions of a string.",
		types.StrType{},
		"a string",
		"repeat",
		[]*params.Param{
			params.SimpleParam("n", "number of repetitions (will be truncated and min'd to 0)", types.NumType{}),
		},
		types.StrType{},
		"the input, repeated n times",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str := string(inputValue.(states.StrValue))
			n := int(argumentValues[0].(states.NumValue))
			if n < 0 {
				n = 0
			}
			return states.StrValue(strings.Repeat(str, n)), nil
		},
		[]shapes.Example{
			{`"abc" repeat(3)`, `Str`, `"abcabcabc"`, nil},
			{`"abc" repeat(0)`, `Str`, `""`, nil},
			{`"abc" repeat(-1)`, `Str`, `""`, nil},
			{`"abc" repeat(1.6)`, `Str`, `"abc"`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Extracts a substring from a string.",
		types.StrType{},
		"a string",
		"slice",
		[]*params.Param{
			params.SimpleParam("start", "a positive integer", types.NumType{}),
		},
		types.StrType{},
		"the portion of the input that is after offset start",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str := string(inputValue.(states.StrValue))
			start := int(argumentValues[0].(states.NumValue))
			if start < 0 {
				start = len(str) + start
				if start < 0 {
					start = 0
				}
			} else if start > len(str) {
				start = len(str)
			}
			return states.StrValue(str[start:]), nil
		},
		[]shapes.Example{
			{`"abc" slice(-4)`, `Str`, `"abc"`, nil},
			{`"abc" slice(-3)`, `Str`, `"abc"`, nil},
			{`"abc" slice(-2)`, `Str`, `"bc"`, nil},
			{`"abc" slice(-1)`, `Str`, `"c"`, nil},
			{`"abc" slice(0)`, `Str`, `"abc"`, nil},
			{`"abc" slice(1)`, `Str`, `"bc"`, nil},
			{`"abc" slice(2)`, `Str`, `"c"`, nil},
			{`"abc" slice(3)`, `Str`, `""`, nil},
			{`"abc" slice(4)`, `Str`, `""`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Extracts a substring from a string.",
		types.StrType{},
		"a string",
		"slice",
		[]*params.Param{
			params.SimpleParam("start", "offset to start after", types.NumType{}),
			params.SimpleParam("end", "offset to end before", types.NumType{}),
		},
		types.StrType{},
		"the portion of the input that is after offset start but before offset end",
		"Negative offsets are counted from the end of the string.",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str := string(inputValue.(states.StrValue))
			start := int(argumentValues[0].(states.NumValue))
			end := int(argumentValues[1].(states.NumValue))
			if start < 0 {
				start = len(str) + start
				if start < 0 {
					start = 0
				}
			} else if start > len(str) {
				start = len(str)
			}
			if end < 0 {
				end = len(str) + end
			}
			if end < start {
				end = start
			} else if end > len(str) {
				end = len(str)
			}
			return states.StrValue(str[start:end]), nil
		},
		[]shapes.Example{
			{`"abc" slice(1, 2)`, `Str`, `"b"`, nil},
			{`"abc" slice(1, -1)`, `Str`, `"b"`, nil},
			{`"abc" slice(-2, -1)`, `Str`, `"b"`, nil},
			{`"abc" slice(-1, -2)`, `Str`, `""`, nil},
			{`"abc" slice(2, 1)`, `Str`, `""`, nil},
			{`"abc" slice(-5, -4)`, `Str`, `""`, nil},
			{`"abc" slice(-5, -3)`, `Str`, `""`, nil},
			{`"abc" slice(-5, -2)`, `Str`, `"a"`, nil},
			{`"abc" slice(-5, -1)`, `Str`, `"ab"`, nil},
			{`"abc" slice(-1, -5)`, `Str`, `""`, nil},
			{`"abc" slice(2, -5)`, `Str`, `""`, nil},
			{`"abc" slice(0, 4)`, `Str`, `"abc"`, nil},
			{`"abc" slice(4, 4)`, `Str`, `""`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Removes whitespace from the start and end of a string.",
		types.StrType{},
		"a string",
		"trim",
		nil,
		types.StrType{},
		"the input, with leading and trailing whitespace removed",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str := string(inputValue.(states.StrValue))
			return states.StrValue(strings.TrimSpace(str)), nil
		},
		[]shapes.Example{
			{`" abc  " trim`, `Str`, `"abc"`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Removes whitespace from the start of a string.",
		types.StrType{},
		"a string",
		"trimStart",
		nil,
		types.StrType{},
		"the input, with leading whitespace removed",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str := string(inputValue.(states.StrValue))
			return states.StrValue(strings.TrimLeftFunc(str, unicode.IsSpace)), nil
		},
		[]shapes.Example{
			{`" abc  " trimStart`, `Str`, `"abc  "`, nil},
		},
	),
	shapes.SimpleFuncer(
		"Removes whitespace from the end of a string.",
		types.StrType{},
		"a string",
		"trimEnd",
		nil,
		types.StrType{},
		"the input, with trailing whitespace removed",
		"",
		func(inputValue states.Value, argumentValues []states.Value) (states.Value, error) {
			str := string(inputValue.(states.StrValue))
			return states.StrValue(strings.TrimRightFunc(str, unicode.IsSpace)), nil
		},
		[]shapes.Example{
			{`" abc  " trimEnd`, `Str`, `" abc"`, nil},
		},
	),
}
