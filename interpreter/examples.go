package interpreter

import (
	"github.com/texttheater/bach/shapes"
)

var SimpleTypeExamples = []shapes.Example{
	{`null`, `Null`, `null`, nil},
	{`false`, `Bool`, `false`, nil},
	{`true`, `Bool`, `true`, nil},
	{`42`, `Num`, `42`, nil},
	{`0.3`, `Num`, `0.3`, nil},
	{`"Hello world!"`, `Str`, `"Hello world!"`, nil},
}

var ArrayTypeExamples = []shapes.Example{
	{`[]`, `Arr<>`, `[]`, nil},
	{`[1]`, `Arr<Num>`, `[1]`, nil},
	{`[1, 2, 3]`, `Arr<Num, Num, Num>`, `[1, 2, 3]`, nil},
	{`[1, 2, 3] each(+1)`, `Arr<Num...>`, `[2, 3, 4]`, nil},
	{`[1, "a"]`, `Arr<Num, Str>`, `[1, "a"]`, nil},
	{`[[1, 2], ["a", "b"]]`, `Arr<Arr<Num, Num>, Arr<Str, Str>>`, `[[1, 2], ["a", "b"]]`, nil},
	{`[1;[]]`, `Arr<Num>`, `[1]`, nil},
	{`[1, 2; [3, 4]]`, `Arr<Num, Num, Num, Num>`, `[1, 2, 3, 4]`, nil},
	{`[3, 4] =rest [1, 2; rest]`, `Arr<Num, Num, Num, Num>`, `[1, 2, 3, 4]`, nil},
	{`[1, 2; [1, 2] each(+2)]`, `Arr<Num, Num, Num...>`, `[1, 2, 3, 4]`, nil},
	{`for Arr<Any...> def f Arr<Any...> as id ok [] f`, `Arr<Any...>`, `[]`, nil},
}

var ObjectTypeExamples = []shapes.Example{
	{`{}`, `Obj<Void>`, `{}`, nil},
	{`{a: 1}`, `Obj<a: Num, Void>`, `{a: 1}`, nil},
	{`{a: 1, b: "c"}`, `Obj<a: Num, b: Str, Void>`, `{a: 1, b: "c"}`, nil},
	{`for Any def f Obj<Num> as {a: 1, b: 2} ok f`, `Obj<Num>`, `{a: 1, b: 2}`, nil},
	{`for Any def f Obj<Any> as {a: 1, b: "c"} ok f`, `Obj<Any>`, `{a: 1, b: "c"}`, nil},
}

var UnionTypeExamples = []shapes.Example{
	{`[1] +["a"]`, `Arr<Num|Str...>`, `[1, "a"]`, nil},
	{`[1] +["a"] get(0)`, `Num|Str`, `1`, nil},
}
