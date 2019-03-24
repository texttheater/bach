package grammar

import (
	"github.com/alecthomas/participle/lexer"
)

var LexerDefinition = lexer.Must(lexer.Regexp(
	`([\s]+)` +
		`|(?P<Num>(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+))` +
		`|(?P<Str>"(?:\\.|[^"])*")` +
		`|(?P<Op1Num>[+\-*/%<>](?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+))` +
		`|(?P<Op2Num>(?:==|<=|>=)(?:\d+\.(?:\d+)?(?:[eE][+-]?\d+)?|\d+[eE][+-]?\d+|\.\d+(?:[eE][+-]?\d+)?|\d+))` +
		`|(?P<Op1Name>[+\-*/%<>](?:[\p{L}_][\p{L}_0-9]*))` +
		`|(?P<Op2Name>(?:==|<=|>=)(?:[\p{L}_][\p{L}_0-9]*))` +
		`|(?P<Assignment>=(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*))` +
		`|(?P<NameLpar>(?:[+\-*/%<>=]|==|<=|>=|[\p{L}_][\p{L}_0-9]*)\()` +
		`|(?P<TypeKeywordLangle>(?:Void|Null|Bool|Num|Str|Seq|Arr|Obj|Any)<)` +
		`|(?P<Prop>[\p{L}_][\p{L}_0-9]*)` +
		`|(?P<Op1>[+\-*/%<>=])` +
		`|(?P<Op2>==|<=|>=)` +
		`|(?P<Comma>,)` +
		`|(?P<Rpar>\))` +
		`|(?P<Lbrack>\[)` +
		`|(?P<Rbrack>])` +
		`|(?P<Lbrace>{)` +
		`|(?P<Rbrace>})` +
		`|(?P<Colon>:)` +
		`|(?P<Pipe>\|)` +
		// the following will be scanned as Name, but mapped to the
		// appropriate token types by ToKeyword (see below)
		`|(?P<Keyword>for|def|as|ok|if|then|elif|else)` +
		`|(?P<TypeKeyword>Void|Null|Bool|Num|Str|Seq|Arr|Obj|Any)`,
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
	return name == "Void" || name == "Null" || name == "Bool" ||
		name == "Num" || name == "Str" || name == "Seq" ||
		name == "Arr" || name == "Any"
}

func isKeyword(name string) bool {
	return name == "for" || name == "def" || name == "as" ||
		name == "ok" || name == "if" || name == "then" ||
		name == "elif" || name == "else" || name == "each" ||
		name == "all"
}