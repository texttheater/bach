package main_test

import (
	"testing"

	"github.com/texttheater/bach/tests"
)

func TestInterp(t *testing.T) {
	tests.TestGrammar(t)
	tests.TestTypeErrors(t)
	tests.TestLiterals(t)
	tests.TestMath(t)
	tests.TestLogic(t)
	tests.TestNull(t)
	tests.TestAssignment(t)
	tests.TestStrings(t)
	tests.TestArrays(t)
	tests.TestObjects(t)
	tests.TestDefinitions(t)
	tests.TestCalls(t)
	tests.TestConditionals(t)
	tests.TestMatchingType(t)
	tests.TestMatchingArr(t)
	tests.TestMatchingObj(t)
	tests.TestRecursion(t)
	tests.TestOverloading(t)
	tests.TestClosures(t)
	tests.TestSimpleTypes(t)
	tests.TestArrayTypes(t)
	tests.TestObjectTypes(t)
	tests.TestUnionTypes(t)
	tests.TestAnyType(t)
	tests.TestVoidType(t)
	tests.TestFilters(t)
	tests.TestValues(t)
	tests.TestRegexp(t)
	tests.TestGetters(t)
}
