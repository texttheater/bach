package grammar

import (
	"github.com/alecthomas/participle/lexer"
)

var LexerDefinition = lexer.Must(lexer.Regexp(
	`([\s]+)` +
		// tokens starting type literals
		`|(?P<TypeKeywordLangle>(?:Void|Null|Bool|Num|Str|Arr|Tup|Obj|Any)<)` +
		// tokens starting calls
		`|(?P<Op1Num>[+\-*/%<>](?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+))` +
		`|(?P<Op2Num>(?:==|<=|>=|\*\*)(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+))` +
		`|(?P<LangleLid><(?:[\p{L}_][\p{L}_0-9]*))` + // special case of Op1Lid, but also used for type variables
		`|(?P<Op1Lid>[+\-*/%>](?:[\p{L}_][\p{L}_0-9]*))` +
		`|(?P<Op2Lid>(?:==|<=|>=|\*\*)(?:[\p{L}_][\p{L}_0-9]*))` +
		`|(?P<NameStr>(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*)"(?:\\.|[^"])*")` +
		`|(?P<NameRegexp>(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*)~(?:\\.|[^/])*)~` +
		`|(?P<NameLpar>(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*)\()` +
		`|(?P<NameLbrack>(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*)\[)` +
		`|(?P<NameLbrace>(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*){)` +
		// assignment
		`|(?P<Assignment>=(?:[+\-*/%<>]|==|<=|>=|[\p{L}_][\p{L}_0-9]*))` +
		// names
		`|(?P<Lid>[\p{L}_][\p{L}_0-9]*)` +
		`|(?P<Op2>==|<=|>=|\*\*)` +
		`|(?P<Op1>[+\-*/%<>=])` +
		// ellipsis
		`|(?P<Ellipsis>\.\.\.)` +
		// numbers
		`|(?P<Num>\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+)` +
		// getters
		`|(?P<LidGetter>@[\p{L}_][\p{L}_0-9]*)` +
		`|(?P<Op1Getter>@[+\-*/%<>=])` +
		`|(?P<Op2Getter>@(?:==|<=|>=))` +
		`|(?P<NumGetter>@(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+))` +
		`|(?P<StrGetter>@"(?:\\.|[^"])*")` +
		// starting with unique characters
		`|(?P<Str>"(?:\\.|[^"])*")` +
		`|(?P<Regexp>~(?:\\.|[^/])*)~` +
		`|(?P<Comma>,)` +
		`|(?P<Lpar>\()` +
		`|(?P<Rpar>\))` +
		`|(?P<Lbrack>\[)` +
		`|(?P<Rbrack>])` +
		`|(?P<Lbrace>{)` +
		`|(?P<Rbrace>})` +
		`|(?P<Colon>:)` +
		`|(?P<Pipe>\|)` +
		`|(?P<Semi>\;)` +
		// the following will be scanned as Name, but mapped to the
		// appropriate token types by ToKeyword (see below)
		`|(?P<Keyword>for|def|as|ok|if|then|elif|else|each|is|elis|with|drop)` +
		`|(?P<TypeKeyword>Void|Null|Reader|Bool|Num|Str|Arr|Tup|Obj|Any)`,
))

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
		name == "with" || name == "drop" || name == "reject"
}
