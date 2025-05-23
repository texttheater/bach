package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/stateful"
)

var LexerDefinition lexer.Definition = lexer.Must(stateful.New(stateful.Rules{
	"Root": {
		// whitespace
		{"whitespace", `\s+`, nil},
		// tokens starting type literals
		{"TypeKeywordLangle", `(?:Arr|Obj)<`, nil},
		{"TypeKeyword", `(?:Void|Null|Reader|Bool|Num|Str|Any)\b`, nil},
		// tokens starting calls
		{"Op1Num", `[+\-*/%<>](?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+)`, nil},
		{"Op2Num", `(?:==|<=|>=|\*\*)(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+)`, nil},
		{"LangleULid", `<(?:[\p{Lu}_][\p{L}_0-9]*)`, nil}, // used for type variables
		{"Op1Lid", `[+\-*/%<>](?:[\p{Ll}_][\p{L}_0-9]*)`, nil},
		{"Op2Lid", `(?:==|<=|>=|\*\*)(?:[\p{Ll}_][\p{L}_0-9]*)`, nil},
		{"NameQuot", `(?:[+\-*/%<>]|==|<=|>=|[\p{Ll}_][\p{L}_0-9]*)"`, stateful.Push("StrLiteral")},
		{"NameRegexp", `(?:[+\-*/%<>]|==|<=|>=|[\p{Ll}_][\p{L}_0-9]*)~(?:\\.|[^/])*~`, nil},
		{"NameLpar", `(?:[+\-*/%<>]|==|<=|>=|\*\*|[\p{Ll}_][\p{L}_0-9]*)\(`, nil},
		{"NameLbrack", `(?:[+\-*/%<>]|==|<=|>=|[\p{Ll}_][\p{L}_0-9]*)\[`, nil},
		{"NameLbrace", `(?:[+\-*/%<>]|==|<=|>=|[\p{Ll}_][\p{L}_0-9]*){`, stateful.Push("Braces")},
		// assignment
		{"EqName", `=(?:[+\-*/%<>]|==|<=|>=|[\p{Ll}_][\p{L}_0-9]*)`, nil},
		{"EqLbrack", `=\[`, nil},
		{"EqLbrace", `={`, stateful.Push("Braces")},
		// keywords
		{"Keyword", `(?:for|def|as|ok|if|then|elif|else|is|elis|with)\b`, nil},
		// names
		{"Lid", `[\p{Ll}_][\p{L}_0-9]*`, nil},
		{"Op2", `==|<=|>=|\*\*`, nil},
		{"Op1", `[+\-*/%<>=]`, nil},
		// ellipsis
		{"Ellipsis", `\.\.\.`, nil},
		// numbers
		{"NumLiteral", `\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+`, nil},
		// getters
		{"NumGetter", `@-?(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+)`, nil},
		{"LidGetter", `@[\p{Ll}_][\p{L}_0-9]*`, nil},
		{"Op2Getter", `@(?:==|<=|>=|\*\*)`, nil},
		{"Op1Getter", `@[+\-*/%<>=]`, nil},
		{"StrGetter", `@"(?:\\.|[^"])*"`, nil},
		// starting with unique characters
		{"Regexp", `~(?:\\.|[^/])*~`, nil},
		{"Comma", `,`, nil},
		{"Lpar", `\(`, nil},
		{"Rpar", `\)`, nil},
		{"Lbrack", `\[`, nil},
		{"Rbrack", `]`, nil},
		{"Lbrace", `{`, stateful.Push("Braces")},
		{"Colon", `:`, nil},
		{"Pipe", `\|`, nil},
		{"Semi", `\;`, nil},
		{"Quot", `"`, stateful.Push("StrLiteral")},
	},
	"StrLiteral": {
		{"Quot", `"`, stateful.Pop()},
		{"Dbrace", `{\{|}}?`, nil},
		{"Lbrace", `{`, stateful.Push("Braces")},
		{"Char", `(?:\\.|[^{}"])+`, nil},
	},
	"Braces": { // generic group for placeholders and object literals
		stateful.Include("Root"),
		{"Rbrace", `}`, stateful.Pop()},
	},
}))
