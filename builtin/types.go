package builtin

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

var TypeFuncers = []shapes.Funcer{
	shapes.Funcer{
		Summary:           "Gives the type of the input expression.",
		InputType:         types.NewVar("A", types.Any{}),
		InputDescription:  "any value (is ignored)",
		Name:              "type",
		Params:            nil,
		OutputType:        types.Str{},
		OutputDescription: "a string representation of the type of the input expression",
		Notes:             "",
		Kernel: func(inputState states.State, args []states.Action, bindings map[string]types.Type, pos lexer.Position) *states.Thunk {
			return states.ThunkFromValue(states.StrValue(bindings["A"].String()))
		},
		IDs: nil,
		Examples: []shapes.Example{
			{`type`, `Str`, `"Null"`, nil},
			{`null type`, `Str`, `"Null"`, nil},
			{`true type`, `Str`, `"Bool"`, nil},
			{`1 type`, `Str`, `"Num"`, nil},
			{`1.5 type`, `Str`, `"Num"`, nil},
			{`"abc" type`, `Str`, `"Str"`, nil},
			{`for Any def f Any as null ok f type`, `Str`, `"Any"`, nil},
			{`[] type`, `Str`, `"Arr<>"`, nil},
			{`["dog", "cat"] type`, `Str`, `"Arr<Str, Str>"`, nil},
			{`["dog", 1] type`, `Str`, `"Arr<Str, Num>"`, nil},
			{`["dog", 1, {}] type`, `Str`, `"Arr<Str, Num, Obj<Void>>"`, nil},
			{`["dog", 1, {}, 2] type`, `Str`, `"Arr<Str, Num, Obj<Void>, Num>"`, nil},
			{`{} type`, `Str`, `"Obj<Void>"`, nil},
			{`{a: null} type`, `Str`, `"Obj<a: Null, Void>"`, nil},
			{`{b: false, a: null} type`, `Str`, `"Obj<a: Null, b: Bool, Void>"`, nil},
			{`{c: 0, b: false, a: null} type`, `Str`, `"Obj<a: Null, b: Bool, c: Num, Void>"`, nil},
			{`{d: "", c: 0, b: false, a: null} type`, `Str`, `"Obj<a: Null, b: Bool, c: Num, d: Str, Void>"`, nil},
			{`{e: [], d: "", c: 0, b: false, a: null} type`, `Str`, `"Obj<a: Null, b: Bool, c: Num, d: Str, e: Arr<>, Void>"`, nil},
			{`{f: {}, e: [], d: "", c: 0, b: false, a: null} type`, `Str`, `"Obj<a: Null, b: Bool, c: Num, d: Str, e: Arr<>, f: Obj<Void>, Void>"`, nil},
			{`for Num def f Num|Str as if ==1 then 1 else "abc" ok ok 1 f type`, `Str`, `"Num|Str"`, nil},
			{`for Any def f Num|Str as 1 ok f type`, `Str`, `"Num|Str"`, nil},
			{`for Any def f Void|Num as 1 ok f type`, `Str`, `"Num"`, nil},
			{`for Any def f Num|Any as 1 ok f type`, `Str`, `"Any"`, nil},
		},
	},
}
