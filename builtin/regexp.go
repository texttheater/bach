package builtin

import (
	"strings"
	"unicode/utf8"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var RegexpFuncers = []shapes.Funcer{
	shapes.Funcer{
		Summary:          "Finds all non-overlapping matches in a string.",
		InputType:        types.Str{},
		InputDescription: "a string",
		Name:             "reFindAll",
		Params: []*params.Param{
			{
				InputType:   types.Str{},
				Name:        "regexp",
				Description: "a [regexp](regexp-expressions.md)",
				Params:      nil,
				OutputType: types.NewVar("A", types.NewUnion(
					types.Null{},
					types.Obj{
						Props: map[string]types.Type{
							"start": types.Num{},
							"0":     types.Str{},
						},
						Rest: types.Any{},
					},
				)),
			},
		},
		OutputType: types.NewArr(
			types.NewVar("A", types.Any{}),
		),
		OutputDescription: "array of matches",
		Notes:             "Matches appear in the output from leftmost to rightmost. Matches that overlap an earlier match (i.e., a match that starts at a lower offset or one that starts at the same offset but is found earlier by the regexp) are not included.",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			str, err := inputState.Thunk.EvalStr()
			if err != nil {
				return states.ThunkFromError(err)
			}
			offset := 0
			iter := func() (*states.Thunk, bool, error) {
				val, err := args[0](inputState.Replace(states.ThunkFromValue(states.StrValue(str))), nil).Thunk.Eval()
				if err != nil {
					return nil, false, err
				}
				objValue, ok := (val.(states.ObjValue))
				if !ok {
					return nil, false, nil
				}
				obj := map[string]*states.Thunk(objValue)
				start, err := obj["start"].EvalInt()
				if err != nil {
					return nil, false, err
				}
				obj["start"] = states.ThunkFromValue(states.NumValue(start + offset))
				group, err := obj["0"].EvalStr()
				if err != nil {
					return nil, false, err
				}
				length := len(group)
				end := start + length
				offset += end
				str = str[end:]
				return states.ThunkFromValue(objValue), true, nil
			}
			return states.ThunkFromIter(iter)
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`"a" reFindAll~a+~`, `Arr<Null|Obj<0: Str, start: Num, Void>...>`, `[{start: 0, 0: "a"}]`, nil},
			{`"aa" reFindAll~a+~`, `Arr<Null|Obj<0: Str, start: Num, Void>...>`, `[{start: 0, 0: "aa"}]`, nil},
			{`"aba" reFindAll~a+~`, `Arr<Null|Obj<0: Str, start: Num, Void>...>`, `[{start: 0, 0: "a"}, {start: 2, 0: "a"}]`, nil},
		},
	},
	shapes.Funcer{
		Summary:          "Replaces the first match of a regexp in a string with something else.",
		InputType:        types.Str{},
		InputDescription: "a string",
		Name:             "reReplaceFirst",
		Params: []*params.Param{
			{
				InputType:   types.Str{},
				Name:        "regexp",
				Description: "a [regexp](regexp-expressions.md)",
				Params:      nil,
				OutputType: types.NewUnion(
					types.Null{},
					types.Obj{
						Props: map[string]types.Type{
							"start": types.Num{},
							"0":     types.Str{},
						},
						Rest: types.Any{},
					},
				),
			},
			{
				InputType: types.Obj{
					Props: map[string]types.Type{
						"start": types.Num{},
						"0":     types.Str{},
					},
					Rest: types.Any{},
				},
				Name:        "replacement",
				Description: "takes a match and returns a string",
				Params:      nil,
				OutputType:  types.Str{},
			},
		},
		OutputType:        types.Str{},
		OutputDescription: "the input with the first match of the regexp replaced with the corresponding replacement, or unchanged if there is no match",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			matchThk := args[0](inputState, nil).Thunk
			match, err := matchThk.Eval()
			if err != nil {
				return states.ThunkFromError(nil)
			}
			switch match := match.(type) {
			case states.NullValue:
				return inputState.Thunk
			case states.ObjValue:
				old, err := inputState.Thunk.EvalStr()
				if err != nil {
					return states.ThunkFromError(nil)
				}
				start, err := match["start"].EvalInt()
				if err != nil {
					return states.ThunkFromError(nil)
				}
				replaced, err := match["0"].EvalStr()
				if err != nil {
					return states.ThunkFromError(nil)
				}
				length := len(replaced)
				replacement, err := args[1](inputState.Replace(matchThk), nil).Thunk.EvalStr()
				if err != nil {
					return states.ThunkFromError(err)
				}
				new_ := old[:start] + replacement + old[start+length:]
				return states.ThunkFromValue(states.StrValue(new_))
			default:
				panic("unexpected type")
			}
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`"aba" reReplaceFirst(~a+~, "c")`, `Str`, `"cba"`, nil},
			{`"abc" reReplaceFirst(~d~, "e")`, `Str`, `"abc"`, nil},
			{`"b0b" reReplaceFirst(~\d+~, @0 parseInt =n "a" repeat(n))`, `Str`, `"bb"`, nil},
			{`"b3b" reReplaceFirst(~\d+~, @0 parseInt =n "a" repeat(n))`, `Str`, `"baaab"`, nil},
			{`" a b c " reReplaceFirst(~[abc]~, "({@0})")`, `Str`, `" (a) b c "`, nil},
		},
	},
	shapes.Funcer{
		Summary:          "Replaces all non-overlapping matches of a regexp in a string with something else.",
		InputType:        types.Str{},
		InputDescription: "a string",
		Name:             "reReplaceAll",
		Params: []*params.Param{
			{
				InputType:   types.Str{},
				Name:        "regexp",
				Description: "a [regexp](regexp-expressions.md)",
				Params:      nil,
				OutputType: types.NewVar("A", types.NewUnion(
					types.Null{},
					types.Obj{
						Props: map[string]types.Type{
							"start": types.Num{},
							"0":     types.Str{},
						},
						Rest: types.Any{},
					},
				)),
			},
			{
				InputType: types.Obj{
					Props: map[string]types.Type{
						"start": types.Num{},
						"0":     types.Str{},
					},
					Rest: types.Any{},
				},
				Name:        "replacement",
				Description: "takes a match and returns a string",
				Params:      nil,
				OutputType:  types.Str{},
			},
		},
		OutputType:        types.Str{},
		OutputDescription: "the input with all matches of the regexp replaced with the corresponding replacement, or unchanged if there is no match",
		Notes:             "Matches are replaced from leftmost to rightmost. Matches that overlap an earlier match (i.e., a match that starts at a lower offset or one that starts at the same offset but is found earlier by the regexp) are not replaced.",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			input, err := inputState.Thunk.EvalStr()
			if err != nil {
				return states.ThunkFromError(err)
			}
			var output strings.Builder
		loop:
			for {
				matchThk := args[0](inputState.Replace(states.ThunkFromValue(states.StrValue(input))), nil).Thunk
				match, err := matchThk.Eval()
				if err != nil {
					return states.ThunkFromError(err)
				}
				switch match := match.(type) {
				case states.NullValue:
					output.WriteString(input)
					break loop
				case states.ObjValue:
					start, err := match["start"].EvalInt()
					if err != nil {
						return states.ThunkFromError(err)
					}
					group, err := match["0"].EvalStr()
					if err != nil {
						return states.ThunkFromError(err)
					}
					length := len(group)
					output.WriteString(input[:start])
					replacement, err := args[1](inputState.Replace(matchThk), nil).Thunk.EvalStr()
					if err != nil {
						return states.ThunkFromError(err)
					}
					output.WriteString(replacement)
					input = input[start+length:]
				}
			}
			return states.ThunkFromValue(states.StrValue(output.String()))
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`"aba" reReplaceAll(~a+~, "c")`, `Str`, `"cbc"`, nil},
			{`"abc" reReplaceAll(~d~, "e")`, `Str`, `"abc"`, nil},
			{`" a b c " reReplaceAll(~[abc]~, "({@0})")`, `Str`, `" (a) (b) (c) "`, nil},
		},
	},
	shapes.Funcer{
		Summary:          "Splits a string around a regexp.",
		InputType:        types.Str{},
		InputDescription: "a string",
		Name:             "reSplit",
		Params: []*params.Param{
			{
				InputType:   types.Str{},
				Name:        "separator",
				Description: "a [regexp](regexp-expressions.md)",
				Params:      nil,
				OutputType: types.NewUnion(
					types.Null{},
					types.Obj{
						Props: map[string]types.Type{
							"start": types.Num{},
							"0":     types.Str{},
						},
						Rest: types.Any{},
					},
				),
			},
		},
		OutputType:        types.NewArr(types.Str{}),
		OutputDescription: "the parts of the input found in between occurrences of the separator",
		Notes:             "If the separator regexp matches the empty string, the input is split into its individual code points. Separators are found from leftmost to rightmost. Separators that overlap an earlier separator (i.e., a separator that starts at a lower offset or one that starts at the same offset but is found earlier by the regexp) do not lead to splits.",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return split(inputState, args, bindings, pos)
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`"zabacad" reSplit~a~`, `Arr<Str...>`, `["z", "b", "c", "d"]`, nil},
			{`"zabaca" reSplit~a~`, `Arr<Str...>`, `["z", "b", "c", ""]`, nil},
			{`"abacad" reSplit~a~`, `Arr<Str...>`, `["", "b", "c", "d"]`, nil},
			{`"abaca" reSplit~a~`, `Arr<Str...>`, `["", "b", "c", ""]`, nil},
			{`"abaca" reSplit~~`, `Arr<Str...>`, `["a", "b", "a", "c", "a"]`, nil},
			{`"你好" reSplit~~`, `Arr<Str...>`, `["你", "好"]`, nil},
			{`"" reSplit~a~`, `Arr<Str...>`, `[""]`, nil},
			{`"" reSplit~~`, `Arr<Str...>`, `[]`, nil},
		},
	},
	shapes.Funcer{
		Summary:          "Splits a string around a regexp, up to a certain number of times.",
		InputType:        types.Str{},
		InputDescription: "a string",
		Name:             "reSplit",
		Params: []*params.Param{
			{
				InputType:   types.Str{},
				Name:        "separator",
				Description: "a [regexp](regexp-expressions.md)",
				Params:      nil,
				OutputType: types.NewUnion(
					types.Null{},
					types.Obj{
						Props: map[string]types.Type{
							"start": types.Num{},
							"0":     types.Str{},
						},
						Rest: types.Any{},
					},
				),
			},
			params.SimpleParam("n", "maximum number of splits to make", types.Num{}),
		}, OutputType: types.NewArr(
			types.Str{},
		),
		OutputDescription: "the parts of the input found in between occurrences of the separator",
		Notes:             "If the separator regexp matches the empty string, the input is split into its individual code points. Separators are found from leftmost to rightmost. Separators that overlap an earlier separator (i.e., a separator that starts at a lower offset or one that starts at the same offset but is found earlier by the regexp) do not lead to splits. At most n splits are made so that the output contains at most n + 1 elements; later separator occurrences are ignored.",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return split(inputState, args, bindings, pos)
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`"zabacad" reSplit(~a~, 1)`, `Arr<Str...>`, `["z", "bacad"]`, nil},
			{`"zabaca" reSplit(~a~, 1)`, `Arr<Str...>`, `["z", "baca"]`, nil},
			{`"zabaca" reSplit(~a~, 3)`, `Arr<Str...>`, `["z", "b", "c", ""]`, nil},
			{`"zabaca" reSplit(~a~, 4)`, `Arr<Str...>`, `["z", "b", "c", ""]`, nil},
			{`"abacad" reSplit(~a~, 1)`, `Arr<Str...>`, `["", "bacad"]`, nil},
			{`"abacad" reSplit(~a~, 2)`, `Arr<Str...>`, `["", "b", "cad"]`, nil},
			{`"abaca" reSplit(~a~, 1)`, `Arr<Str...>`, `["", "baca"]`, nil},
			{`"abaca" reSplit(~a~, 2)`, `Arr<Str...>`, `["", "b", "ca"]`, nil},
			{`"abaca" reSplit(~a~, 3)`, `Arr<Str...>`, `["", "b", "c", ""]`, nil},
			{`"abaca" reSplit(~~, 2)`, `Arr<Str...>`, `["a", "b", "aca"]`, nil},
			{`"abaca" reSplit(~~, 1000)`, `Arr<Str...>`, `["a", "b", "a", "c", "a"]`, nil},
			{`"你好" reSplit(~~, 0)`, `Arr<Str...>`, `["你好"]`, nil},
			{`"" reSplit(~a~, 0)`, `Arr<Str...>`, `[""]`, nil},
			{`"" reSplit(~a~, 1)`, `Arr<Str...>`, `[""]`, nil},
			{`"" reSplit(~~, 0)`, `Arr<Str...>`, `[]`, nil},
			{`"" reSplit(~~, 1)`, `Arr<Str...>`, `[]`, nil},
		},
	},
}

func split(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
	str, err := inputState.Thunk.EvalStr()
	if err != nil {
		return states.ThunkFromError(err)
	}
	regexp := args[0]
	maxSplits := -1
	if len(args) > 1 {
		var err error
		maxSplits, err = args[1](inputState.Clear(), nil).Thunk.EvalInt()
		if err != nil {
			return states.ThunkFromError(err)
		}
	}
	splits := 0
	// Edge case 1: if the separator is empty, split the string into
	// codepoints.
	val, err := regexp(inputState.Replace(states.ThunkFromValue(states.StrValue(""))), nil).Thunk.Eval()
	if err != nil {
		return states.ThunkFromError(err)
	}
	_, ok := val.(states.ObjValue)
	if ok {
		iter := func() (*states.Thunk, bool, error) {
			if len(str) == 0 {
				return nil, false, nil
			}
			var l int
			if splits == maxSplits {
				piece := str
				str = ""
				return states.ThunkFromValue(states.StrValue(piece)), true, nil
			}
			_, l = utf8.DecodeRuneInString(str)
			piece := str[:l]
			str = str[l:]
			splits++
			return states.ThunkFromValue(states.StrValue(piece)), true, nil
		}
		return states.ThunkFromIter(iter)
	}
	// Edge case 2: if the string is empty, return a single-element
	// list containing the empty string.
	if str == "" {
		return states.ThunkFromSlice([]*states.Thunk{
			states.ThunkFromValue(states.StrValue("")),
		})
	}
	// Now for the normal cases.
	sepAtEnd := false
	iter := func() (*states.Thunk, bool, error) {
		if sepAtEnd {
			sepAtEnd = false
			return states.ThunkFromValue(states.StrValue("")), true, nil
		}
		if len(str) == 0 {
			// end of string reached
			return nil, false, nil
		}
		if splits == maxSplits {
			piece := str
			str = ""
			return states.ThunkFromValue(states.StrValue(piece)), true, nil
		}
		val, err := regexp(inputState.Replace(states.ThunkFromValue(states.StrValue(str))), nil).Thunk.Eval()
		if err != nil {
			return nil, false, err
		}
		objValue, ok := val.(states.ObjValue)
		if !ok {
			piece := str
			str = ""
			return states.ThunkFromValue(states.StrValue(piece)), true, nil
		}
		obj := map[string]*states.Thunk(objValue)
		sepStart, err := obj["start"].EvalInt()
		if err != nil {
			return nil, false, err
		}
		sep, err := obj["0"].EvalStr()
		if err != nil {
			return nil, false, err
		}
		sepLength := len(sep)
		sepEnd := sepStart + sepLength
		piece := str[:sepStart]
		str = str[sepEnd:]
		splits++
		if len(str) == 0 && sepLength > 0 {
			sepAtEnd = true
		}
		return states.ThunkFromValue(states.StrValue(piece)), true, nil
	}
	return states.ThunkFromIter(iter)
}
