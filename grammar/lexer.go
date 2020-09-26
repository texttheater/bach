package grammar

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/stateful"
)

var LexerDefinition lexer.Definition = lexer.Must(stateful.New(stateful.Rules{
	"Root": {
		{"whitespace", `\s+`, nil},
		// tokens starting type literals
		{"TypeKeywordLangle", `Arr<|Tup<|Obj<`, nil},
		// tokens starting calls
		{"Op1Num", `[+\-*/%<>](?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+)`, nil},
		{"Op2Num", `(?:==|<=|>=|\*\*)(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+)`, nil},
		{"LangleLid", `<(?:[\p{L}_][\p{L}_0-9]*)`, nil}, // special case of Op1Lid, but also used for type variables
		{"Op1Lid", `[+\-*/%>](?:[\p{L}_][\p{L}_0-9]*)`, nil},
		{"Op2Lid", `(?:==|<=|>=|\*\*)(?:[\p{L}_][\p{L}_0-9]*)`, nil},
		{"NameStr", `(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*)"(?:\\.|[^"])*"`, nil},
		{"NameRegexp", `(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*)~(?:\\.|[^/])*~`, nil},
		{"NameLpar", `(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*)\(`, nil},
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
		{"Num", `\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+`, nil},
		// getters
		{"LidGetter", `@[\p{L}_][\p{L}_0-9]*`, nil},
		{"Op1Getter", `@[+\-*/%<>=]`, nil},
		{"Op2Getter", `@(?:==|<=|>=)`, nil},
		{"NumGetter", `@(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+)`, nil},
		{"StrGetter", `@"(?:\\.|[^"])*"`, nil},
		// starting with unique characters
		{"Str", `"(?:\\.|[^"])*"`, nil},
		{"Regexp", `~(?:\\.|[^/])*~`, nil},
		{"Comma", `,`, nil},
		{"Lpar", `\(`, nil},
		{"Rpar", `\)`, nil},
		{"Lbrack", `\[`, nil},
		{"Rbrack", `]`, nil},
		{"Lbrace", `{`, nil},
		{"Rbrace", `}`, nil},
		{"Colon", `:`, nil},
		{"Pipe", `\|`, nil},
		{"Semi", `\;`, nil},
		// the following will be scanned as Name, but mapped to the
		// appropriate token types by ToKeyword (see below)
		{"Keyword", `for|def|as|ok|if|then|elif|else|each|is|elis|with|drop`, nil},
		{"TypeKeyword", `Void|Null|Reader|Bool|Num|Str|Arr|Tup|Obj|Any`, nil},
	},
}))

func ToKeyword(t lexer.Token) (lexer.Token, error) {
	if isTypeKeyword(t.Value) {
		t.Type = LexerDefinition.Symbols()["Type"]
	}
	if isKeyword(t.Value) {
		t.Type = LexerDefinition.Symbols()["Keyword"]
	}
	return t, nil
}

func isTypeKeyword(name string) bool {
	return name == "Void" || name == "Null" || name == "Reader" ||
		name == "Bool" || name == "Num" || name == "Str" ||
		name == "Arr" || name == "Tup" || name == "Obj" ||
		name == "Any"
}

func isKeyword(name string) bool {
	return name == "for" || name == "def" || name == "as" ||
		name == "ok" || name == "if" || name == "then" ||
		name == "elif" || name == "else" || name == "each" ||
		name == "all" || name == "is" || name == "elis" ||
		name == "with" || name == "drop"
}
