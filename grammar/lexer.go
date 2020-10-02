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
		{"Keyword", `(?:for|def|as|ok|if|then|elif|else|each|all|is|elis|with|drop)\b`, nil},
		{"TypeKeywordLangle", `(?:Arr|Tup|Obj)<`, nil},
		{"TypeKeyword", `(?:Void|Null|Reader|Bool|Num|Str|Arr|Tup|Obj|Any)\b`, nil},
		// tokens starting calls
		{"Op1Num", `[+\-*/%<>](?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+)`, nil},
		{"Op2Num", `(?:==|<=|>=|\*\*)(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+)`, nil},
		{"LangleLid", `<(?:[\p{L}_][\p{L}_0-9]*)`, nil}, // special case of Op1Lid, but also used for type variables
		{"Op1Lid", `[+\-*/%>](?:[\p{L}_][\p{L}_0-9]*)`, nil},
		{"Op2Lid", `(?:==|<=|>=|\*\*)(?:[\p{L}_][\p{L}_0-9]*)`, nil},
		{"NameStr", `(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*)"(?:\\.|[^"])*"`, nil},
		{"NameRegexp", `(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*)~(?:\\.|[^/])*~`, nil},
		{"NameLpar", `(?:[+\-*/%<>=]|==|<=|>=|\*\*|[\p{L}_][\p{L}_0-9]*)\(`, nil},
		{"NameLbrack", `(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*)\[`, nil},
		{"NameLbrace", `(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*){`, nil},
		// assignment
		{"Assignment", `=(?:[+\-*/%<>]|==|<=|>=|[\p{L}_][\p{L}_0-9]*)`, nil},
		// names
		{"Lid", `[\p{L}_][\p{L}_0-9]*`, nil},
		{"Op2", `==|<=|>=|\*\*`, nil},
		{"Op1", `[+\-*/%<>=]`, nil},
		// ellipsis
		{"Ellipsis", `\.\.\.`, nil},
		// numbers
		{"NumLiteral", `\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+`, nil},
		// getters
		{"LidGetter", `@[\p{L}_][\p{L}_0-9]*`, nil},
		{"Op1Getter", `@[+\-*/%<>=]`, nil},
		{"Op2Getter", `@(?:==|<=|>=)`, nil},
		{"NumGetter", `@(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+)`, nil},
		{"StrGetter", `@"(?:\\.|[^"])*"`, nil},
		// starting with unique characters
		{"StrLiteral", `"(?:\\.|[^"])*"`, nil},
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
		{"Backtick", "`", stateful.Push("TemplateLiteral")},
	},
	"TemplateLiteral": {
		{"Backtick", "`", stateful.Pop()},
		{"Lbrace", `{`, stateful.Push("Braces")},
		{"Char", "[^{`]+", nil},
	},
	"Braces": { // generic group for placeholders and object literals
		stateful.Include("Root"),
		{"Rbrace", `}`, stateful.Pop()},
	},
}))
