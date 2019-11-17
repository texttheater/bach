package tests

import (
	//"log"
	"reflect"
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

func TestProgram(program string, wantType types.Type, wantValue values.Value, wantError error, t *testing.T) {
	//log.Print(program)
	gotType, gotValue, gotErr := interpreter.InterpretString(program)
	if wantError != nil {
		if gotErr == nil {
			t.Log("ERROR: Expected error but program succeeded.")
			t.Logf("Program:        %s", program)
			t.Logf("Expected error: %s", wantError)
			t.Logf("Got type:       %s", gotType)
			t.Logf("Got value:      %s", gotValue)
			t.Fail()
		} else if !errors.Match(wantError, gotErr) {
			t.Log("ERROR: Expected error does not match actual error.")
			t.Logf("Program:        %s", program)
			t.Logf("Expected error: %s", wantError)
			t.Logf("Got error:      %s", gotErr)
			t.Fail()
		}
	} else {
		if gotErr != nil {
			t.Log("ERROR: Expected program to succeed but got error.")
			t.Logf("Program:        %s", program)
			t.Logf("Expected type:  %s", wantType)
			t.Logf("Expected value: %s", wantValue)
			t.Logf("Got error:      %s", gotErr)
			t.Fail()
		} else if !reflect.DeepEqual(wantType, gotType) {
			t.Log("ERROR: Program has unexpected output type.")
			t.Logf("Program:        %s", program)
			t.Logf("Expected type:  %s", wantType)
			t.Logf("Expected value: %s", wantValue)
			t.Logf("Got type:       %s", gotType)
			t.Logf("Got value:      %s", gotValue)
			t.Fail()
		} else if !wantValue.Equal(gotValue) {
			t.Log("ERROR: Program has unexpected output value.")
			t.Logf("Program:        %s", program)
			t.Logf("Expected type:  %s", wantType)
			t.Logf("Expected value: %s", wantValue)
			t.Logf("Got type:       %s", gotType)
			t.Logf("Got value:      %s", gotValue)
			t.Fail()
		}
	}
}
