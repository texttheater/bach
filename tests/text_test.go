package tests_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/tests"
	"github.com/texttheater/bach/types"
)

func TestText(t *testing.T) {
	tests.TestProgram(
		`"abc"`,
		types.Str{},
		states.StrValue("abc"),
		nil,
		t,
	)
	tests.TestProgram(
		`"\"\\abc\""`,
		types.Str{},
		states.StrValue(`"\abc"`),
		nil,
		t,
	)
	tests.TestProgram(
		`1 "abc"`,
		types.Str{},
		states.StrValue("abc"),
		nil,
		t,
	)
	tests.TestProgram(
		`"a" bytes len`,
		types.Num{},
		states.NumValue(1),
		nil,
		t,
	)
	tests.TestProgram(
		`"ü" bytes len`,
		types.Num{},
		states.NumValue(2),
		nil,
		t,
	)
	tests.TestProgram(
		`"a" codePoints len`,
		types.Num{},
		states.NumValue(1),
		nil,
		t,
	)
	tests.TestProgram(
		`"ü" codePoints len`,
		types.Num{},
		states.NumValue(1),
		nil,
		t,
	)
	tests.TestProgram(
		`"a{2 +2}"`,
		types.Str{},
		states.StrValue("a4"),
		nil,
		t,
	)
	tests.TestProgram(
		`"{{}}"`,
		types.Str{},
		states.StrValue("{}"),
		nil,
		t,
	)
	tests.TestProgram(
		`"{{}"`,
		nil,
		nil,
		errors.SyntaxError(
			errors.Code(errors.Syntax),
		),
		t,
	)
	tests.TestProgram(
		`"\t"`,
		types.Str{},
		states.StrValue("\t"),
		nil,
		t,
	)
	tests.TestProgram(
		`"日本語"`,
		types.Str{},
		states.StrValue("日本語"),
		nil,
		t,
	)
	tests.TestProgram(
		`"\u65e5\u672c\u8a9e"`,
		types.Str{},
		states.StrValue("日本語"),
		nil,
		t,
	)
	tests.TestProgram(
		`"\U000065e5\U0000672c\U00008a9e"`,
		types.Str{},
		states.StrValue("日本語"),
		nil,
		t,
	)
	tests.TestProgram(
		`"\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e"`,
		types.Str{},
		states.StrValue("日本語"),
		nil,
		t,
	)
	tests.TestProgramStr(
		`"  foo bar  baz   " fields`,
		`Arr<Str>`,
		`["foo", "bar", "baz"]`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" startsWith"ab"`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" startsWith"b"`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" endsWith"bc"`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" endsWith"b"`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" slice(0)`,
		`Str`,
		`"abc"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" slice(1)`,
		`Str`,
		`"bc"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" slice(1, 2)`,
		`Str`,
		`"b"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" slice(1, -1)`,
		`Str`,
		`"b"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" slice(-2, -1)`,
		`Str`,
		`"b"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" slice(-1, -2)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" slice(2, 1)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" slice(-2)`,
		`Str`,
		`"bc"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" slice(-4)`,
		`Str`,
		`"abc"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" slice(-5, -4)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" slice(-5, -3)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" slice(-5, -2)`,
		`Str`,
		`"a"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" slice(-5, -1)`,
		`Str`,
		`"ab"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" slice(-1, -5)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" slice(2, -5)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" repeat(3)`,
		`Str`,
		`"abcabcabc"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" repeat(0)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" repeat(-1)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" repeat(1.6)`,
		`Str`,
		`"abc"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" indexOf"b"`,
		`Num`,
		`1`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"abc" indexOf"d"`,
		`Num`,
		`-1`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"ababa" replaceFirst("b", "c")`,
		`Str`,
		`"acaba"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"ababa" replaceAll("b", "c")`,
		`Str`,
		`"acaca"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"z" padEnd(2, " ")`,
		`Str`,
		`"z "`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"z" padEnd(3, " ")`,
		`Str`,
		`"z  "`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"zzz" padEnd(3, " ")`,
		`Str`,
		`"zzz"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"zzzz" padEnd(3, " ")`,
		`Str`,
		`"zzzz"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"zzzz" padEnd(7, "ab")`,
		`Str`,
		`"zzzzaba"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"z" padStart(2, " ")`,
		`Str`,
		`" z"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"z" padStart(3, " ")`,
		`Str`,
		`"  z"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"zzz" padStart(3, " ")`,
		`Str`,
		`"zzz"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"zzzz" padStart(3, " ")`,
		`Str`,
		`"zzzz"`,
		nil,
		t,
	)
	tests.TestProgramStr(
		`"zzzz" padStart(7, "ab")`,
		`Str`,
		`"abazzzz"`,
		nil,
		t,
	)
}
