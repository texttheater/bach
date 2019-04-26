package tests

import (
	"reflect"
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type TestCase struct {
	Program   string
	WantType  types.Type
	WantValue values.Value
	WantError error
}

func (c TestCase) Run(t *testing.T) {
	gotType, gotValue, gotErr := interpreter.InterpretString(c.Program)
	if c.WantError != nil {
		if gotErr == nil {
			t.Log("ERROR: Expected error but program succeeded.")
			t.Logf("Program:        %s", c.Program)
			t.Logf("Expected error: %s", c.WantError)
			t.Logf("Got type:       %s", gotType)
			t.Logf("Got value:      %s", gotValue)
			t.Fail()
		} else if !errors.Match(c.WantError, gotErr) {
			t.Log("ERROR: Expected error does not match actual error.")
			t.Logf("Program:        %s", c.Program)
			t.Logf("Expected error: %s", c.WantError)
			t.Logf("Got error:      %s", gotErr)
			t.Fail()
		}
	} else {
		if gotErr != nil {
			t.Log("ERROR: Expected program to succeed but got error.")
			t.Logf("Program:        %s", c.Program)
			t.Logf("Expected type:  %s", c.WantType)
			t.Logf("Expected value: %s", c.WantValue)
			t.Logf("Got error:      %s", gotErr)
			t.Fail()
		} else if !reflect.DeepEqual(c.WantType, gotType) {
			t.Log("ERROR: Program has unexpected output type.")
			t.Logf("Program:        %s", c.Program)
			t.Logf("Expected type:  %s", c.WantType)
			t.Logf("Expected value: %s", c.WantValue)
			t.Logf("Got type:       %s", gotType)
			t.Logf("Got value:      %s", gotValue)
			t.Fail()
		} else if !reflect.DeepEqual(c.WantValue, gotValue) {
			t.Log("ERROR: Program has unexpected output value.")
			t.Logf("Program:        %s", c.Program)
			t.Logf("Expected type:  %s", c.WantType)
			t.Logf("Expected value: %s", c.WantValue)
			t.Logf("Got type:       %s", gotType)
			t.Logf("Got value:      %s", gotValue)
			t.Fail()
		}
	}
}

func Run(cases []TestCase, t *testing.T) {
	for _, c := range cases {
		c.Run(t)
	}
}
