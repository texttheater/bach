package interpreter

import (
	"github.com/texttheater/bach/shapes"
)

var SimpleTypeExamples = []shapes.Example{
	{`null type`, `Str`, `"Null"`, nil},
	{`false type`, `Str`, `"Bool"`, nil},
	{`true type`, `Str`, `"Bool"`, nil},
	{`42 type`, `Str`, `"Num"`, nil},
	{`0.3 type`, `Str`, `"Num"`, nil},
	{`"Hello world!" type`, `Str`, `"Str"`, nil},
}
