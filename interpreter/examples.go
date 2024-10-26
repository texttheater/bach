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
