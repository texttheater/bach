package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestStrings(t *testing.T) {
	TestProgram(
		`"abc"`,
		types.StrType{},
		states.StrValue("abc"),
		nil,
		t,
	)
	TestProgram(
		`"\"\\abc\""`,
		types.StrType{},
		states.StrValue(`"\abc"`),
		nil,
		t,
	)
	TestProgram(
		`1 "abc"`,
		types.StrType{},
		states.StrValue("abc"),
		nil,
		t,
	)
	TestProgram(
		`"a" bytes len`,
		types.NumType{},
		states.NumValue(1),
		nil,
		t,
	)
	TestProgram(
		`"ü" bytes len`,
		types.NumType{},
		states.NumValue(2),
		nil,
		t,
	)
	TestProgram(
		`"a{2 +2}"`,
		types.StrType{},
		states.StrValue("a4"),
		nil,
		t,
	)
	TestProgram(
		`"{{}}"`,
		types.StrType{},
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
		types.StrType{},
		states.StrValue("\t"),
		nil,
		t,
	)
	TestProgram(
		`"日本語"`,
		types.StrType{},
		states.StrValue("日本語"),
		nil,
		t,
	)
	TestProgram(
		`"\u65e5\u672c\u8a9e"`,
		types.StrType{},
		states.StrValue("日本語"),
		nil,
		t,
	)
	TestProgram(
		`"\U000065e5\U0000672c\U00008a9e"`,
		types.StrType{},
		states.StrValue("日本語"),
		nil,
		t,
	)
	TestProgram(
		`"\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e"`,
		types.StrType{},
		states.StrValue("日本語"),
		nil,
		t,
	)
}
