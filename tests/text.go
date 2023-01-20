package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestText(t *testing.T) {
	TestProgram(
		`"abc"`,
		types.Str{},
		states.StrValue("abc"),
		nil,
		t,
	)
	TestProgram(
		`"\"\\abc\""`,
		types.Str{},
		states.StrValue(`"\abc"`),
		nil,
		t,
	)
	TestProgram(
		`1 "abc"`,
		types.Str{},
		states.StrValue("abc"),
		nil,
		t,
	)
	TestProgram(
		`"a" bytes len`,
		types.Num{},
		states.NumValue(1),
		nil,
		t,
	)
	TestProgram(
		`"ü" bytes len`,
		types.Num{},
		states.NumValue(2),
		nil,
		t,
	)
	TestProgram(
		`"a" codePoints len`,
		types.Num{},
		states.NumValue(1),
		nil,
		t,
	)
	TestProgram(
		`"ü" codePoints len`,
		types.Num{},
		states.NumValue(1),
		nil,
		t,
	)
	TestProgram(
		`"a{2 +2}"`,
		types.Str{},
		states.StrValue("a4"),
		nil,
		t,
	)
	TestProgram(
		`"{{}}"`,
		types.Str{},
		states.StrValue("{}"),
		nil,
		t,
	)
	TestProgram(
		`"{{}"`,
		nil,
		nil,
		errors.SyntaxError(
			errors.Code(errors.Syntax),
		),
		t,
	)
	TestProgram(
		`"\t"`,
		types.Str{},
		states.StrValue("\t"),
		nil,
		t,
	)
	TestProgram(
		`"日本語"`,
		types.Str{},
		states.StrValue("日本語"),
		nil,
		t,
	)
	TestProgram(
		`"\u65e5\u672c\u8a9e"`,
		types.Str{},
		states.StrValue("日本語"),
		nil,
		t,
	)
	TestProgram(
		`"\U000065e5\U0000672c\U00008a9e"`,
		types.Str{},
		states.StrValue("日本語"),
		nil,
		t,
	)
	TestProgram(
		`"\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e"`,
		types.Str{},
		states.StrValue("日本語"),
		nil,
		t,
	)
	TestProgramStr(
		`"  foo bar  baz   " fields`,
		`Arr<Str>`,
		`["foo", "bar", "baz"]`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" startsWith"ab"`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" startsWith"b"`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" endsWith"bc"`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" endsWith"b"`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" slice(0)`,
		`Str`,
		`"abc"`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" slice(1)`,
		`Str`,
		`"bc"`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" slice(1, 2)`,
		`Str`,
		`"b"`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" slice(1, -1)`,
		`Str`,
		`"b"`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" slice(-2, -1)`,
		`Str`,
		`"b"`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" slice(-1, -2)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" slice(2, 1)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" slice(-2)`,
		`Str`,
		`"bc"`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" slice(-4)`,
		`Str`,
		`"abc"`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" slice(-5, -4)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" slice(-5, -3)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" slice(-5, -2)`,
		`Str`,
		`"a"`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" slice(-5, -1)`,
		`Str`,
		`"ab"`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" slice(-1, -5)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" slice(2, -5)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" repeat(3)`,
		`Str`,
		`"abcabcabc"`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" repeat(0)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" repeat(-1)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" repeat(1.6)`,
		`Str`,
		`"abc"`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" indexOf"b"`,
		`Num`,
		`1`,
		nil,
		t,
	)
	TestProgramStr(
		`"abc" indexOf"d"`,
		`Num`,
		`-1`,
		nil,
		t,
	)
	TestProgramStr(
		`"ababa" replaceFirst("b", "c")`,
		`Str`,
		`"acaba"`,
		nil,
		t,
	)
	TestProgramStr(
		`"ababa" replaceAll("b", "c")`,
		`Str`,
		`"acaca"`,
		nil,
		t,
	)
}
