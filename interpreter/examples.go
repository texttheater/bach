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

var TupleTypeExamples = []shapes.Example{
	{`[]`, `Tup<>`, `[]`, nil},
	{`[1]`, `Tup<Num>`, `[1]`, nil},
	{`[1, 2, 3]`, `Tup<Num, Num, Num>`, `[1, 2, 3]`, nil},
	{`[1, "a"]`, `Tup<Num, Str>`, `[1, "a"]`, nil},
	{`[[1, 2], ["a", "b"]]`, `Tup<Tup<Num, Num>, Tup<Str, Str>>`, `[[1, 2], ["a", "b"]]`, nil},
	{`[1;[]]`, `Tup<Num>`, `[1]`, nil},
	{`[1, 2; [3, 4]]`, `Tup<Num, Num, Num, Num>`, `[1, 2, 3, 4]`, nil},
	{`[3, 4] =rest [1, 2; rest]`, `Tup<Num, Num, Num, Num>`, `[1, 2, 3, 4]`, nil},
	{`[1, 2; [1, 2] each(+2)]`, `Tup<Num, Num, Num...>`, `[1, 2, 3, 4]`, nil},
}
