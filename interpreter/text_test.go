package interpreter_test

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestText(t *testing.T) {
	interpreter.TestProgram(
		`"abc"`,
		types.Str{},
		states.StrValue("abc"),
		nil,
		t,
	)
	interpreter.TestProgram(
		`"\"\\abc\""`,
		types.Str{},
		states.StrValue(`"\abc"`),
		nil,
		t,
	)
	interpreter.TestProgram(
		`1 "abc"`,
		types.Str{},
		states.StrValue("abc"),
		nil,
		t,
	)
	interpreter.TestProgram(
		`"a" bytes len`,
		types.Num{},
		states.NumValue(1),
		nil,
		t,
	)
	interpreter.TestProgram(
		`"ü" bytes len`,
		types.Num{},
		states.NumValue(2),
		nil,
		t,
	)
	interpreter.TestProgram(
		`"a" codePoints len`,
		types.Num{},
		states.NumValue(1),
		nil,
		t,
	)
	interpreter.TestProgram(
		`"ü" codePoints len`,
		types.Num{},
		states.NumValue(1),
		nil,
		t,
	)
	interpreter.TestProgram(
		`"a{2 +2}"`,
		types.Str{},
		states.StrValue("a4"),
		nil,
		t,
	)
	interpreter.TestProgram(
		`"{{}}"`,
		types.Str{},
		states.StrValue("{}"),
		nil,
		t,
	)
	interpreter.TestProgram(
		`"{{}"`,
		nil,
		nil,
		errors.SyntaxError(
			errors.Code(errors.Syntax),
		),
		t,
	)
	interpreter.TestProgram(
		`"\t"`,
		types.Str{},
		states.StrValue("\t"),
		nil,
		t,
	)
	interpreter.TestProgram(
		`"日本語"`,
		types.Str{},
		states.StrValue("日本語"),
		nil,
		t,
	)
	interpreter.TestProgram(
		`"\u65e5\u672c\u8a9e"`,
		types.Str{},
		states.StrValue("日本語"),
		nil,
		t,
	)
	interpreter.TestProgram(
		`"\U000065e5\U0000672c\U00008a9e"`,
		types.Str{},
		states.StrValue("日本語"),
		nil,
		t,
	)
	interpreter.TestProgram(
		`"\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e"`,
		types.Str{},
		states.StrValue("日本語"),
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"  foo bar  baz   " fields`,
		`Arr<Str>`,
		`["foo", "bar", "baz"]`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" startsWith"ab"`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" startsWith"b"`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" endsWith"bc"`,
		`Bool`,
		`true`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" endsWith"b"`,
		`Bool`,
		`false`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" slice(0)`,
		`Str`,
		`"abc"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" slice(1)`,
		`Str`,
		`"bc"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" slice(1, 2)`,
		`Str`,
		`"b"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" slice(1, -1)`,
		`Str`,
		`"b"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" slice(-2, -1)`,
		`Str`,
		`"b"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" slice(-1, -2)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" slice(2, 1)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" slice(-2)`,
		`Str`,
		`"bc"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" slice(-4)`,
		`Str`,
		`"abc"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" slice(-5, -4)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" slice(-5, -3)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" slice(-5, -2)`,
		`Str`,
		`"a"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" slice(-5, -1)`,
		`Str`,
		`"ab"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" slice(-1, -5)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" slice(2, -5)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" repeat(3)`,
		`Str`,
		`"abcabcabc"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" repeat(0)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" repeat(-1)`,
		`Str`,
		`""`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" repeat(1.6)`,
		`Str`,
		`"abc"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" indexOf"b"`,
		`Num`,
		`1`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"abc" indexOf"d"`,
		`Num`,
		`-1`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"ababa" replaceFirst("b", "c")`,
		`Str`,
		`"acaba"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"ababa" replaceAll("b", "c")`,
		`Str`,
		`"acaca"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"z" padEnd(2, " ")`,
		`Str`,
		`"z "`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"z" padEnd(3, " ")`,
		`Str`,
		`"z  "`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"zzz" padEnd(3, " ")`,
		`Str`,
		`"zzz"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"zzzz" padEnd(3, " ")`,
		`Str`,
		`"zzzz"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"zzzz" padEnd(7, "ab")`,
		`Str`,
		`"zzzzaba"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"z" padStart(2, " ")`,
		`Str`,
		`" z"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"z" padStart(3, " ")`,
		`Str`,
		`"  z"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"zzz" padStart(3, " ")`,
		`Str`,
		`"zzz"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"zzzz" padStart(3, " ")`,
		`Str`,
		`"zzzz"`,
		nil,
		t,
	)
	interpreter.TestProgramStr(
		`"zzzz" padStart(7, "ab")`,
		`Str`,
		`"abazzzz"`,
		nil,
		t,
	)
}
