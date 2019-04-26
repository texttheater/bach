package main_test

import (
	"testing"

	"github.com/texttheater/bach/tests"
)

func TestInterp(t *testing.T) {
	tests.Run(tests.SyntaxErrorTestCases(), t)
	tests.Run(tests.TypeErrorTestCases(), t)
	tests.Run(tests.LiteralTestCases(), t)
	tests.Run(tests.MathTestCases(), t)
	tests.Run(tests.LogicTestCases(), t)
	tests.Run(tests.NullTestCases(), t)
	tests.Run(tests.AssignmentTestCases(), t)
	tests.Run(tests.StringTestCases(), t)
	tests.Run(tests.ArrayTestCases(), t)
	tests.Run(tests.DefinitionTestCases(), t)
	tests.Run(tests.CallTestCases(), t)
	tests.Run(tests.ConditionalTestCases(), t)
	tests.Run(tests.RecursionTestCases(), t)
	tests.Run(tests.OverloadingTestCases(), t)
	tests.Run(tests.ClosureTestCases(), t)
	tests.Run(tests.SequenceTestCases(), t)
	tests.Run(tests.SimpleTypeTestCases(), t)
	tests.Run(tests.SequenceTypeTestCases(), t)
	tests.Run(tests.ArrayTypeTestCases(), t)
	tests.Run(tests.ObjectTypeTestCases(), t)
	tests.Run(tests.UnionTypeTestCases(), t)
	tests.Run(tests.AnyTypeTestCases(), t)
}
