package tests

import (
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestMatchingArr(t *testing.T) {
	TestProgram(
		`[1, 2, 3] is [Num, Num, Num] then true ok`,
		types.BoolType{},
		values.BoolValue(true),
		nil,
		t)
	TestProgram(
		`[1, 2, 3] is [Num, Num, Num] then true else false ok`,
		nil,
		nil,
		errors.E(
			errors.Code(errors.UnreachableElseClause),
		),
		t)
	TestProgram(
		`[1, 2, 3] each id all arr is [Num, Num, Num] then true else false ok`,
		types.BoolType{},
		values.BoolValue(true),
		nil,
		t)
	TestProgram(
		`[1, "a"] is [Num, Str] then true ok`,
		types.BoolType{},
		values.BoolValue(true),
		nil,
		t)
	TestProgram(
		`[1, "a"] is [Num a, Str b] then a ok`,
		types.NumType{},
		values.NumValue(1),
		nil,
		t)
	TestProgram(
		`[1, "a"] is [Num a, Str b] then b ok`,
		types.StrType{},
		values.StrValue("a"),
		nil,
		t)
	TestProgram(
		`[[1]] is [[Any x]] then x ok`,
		types.NumType{},
		values.NumValue(1),
		nil,
		t)
}
