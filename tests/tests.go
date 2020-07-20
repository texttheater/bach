package tests

import (
	//"log"
	"testing"

	"github.com/texttheater/bach/interpreter"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestProgram(program string, wantType types.Type, wantValue states.Value, wantError error, t *testing.T) {
	//log.Print(program)
	var wantValueStr string
	if wantValue != nil {
		wantValueStr, _ = wantValue.String()
	}
	gotType, gotValue, gotErr := interpreter.InterpretString(program)
	var gotValueStr string
	if gotValue != nil {
		var err error
		gotValueStr, err = gotValue.String()
		if gotErr == nil {
			gotErr = err
		}
	}
	if wantError != nil {
		if gotErr == nil {
			t.Log("ERROR: Expected error but program succeeded.")
			t.Logf("Program:        %s", program)
			t.Logf("Expected error: %s", wantError)
			t.Logf("Got type:       %s", gotType)
			t.Logf("Got value:      %s", gotValueStr)
			t.Fail()
		} else if !states.Match(wantError, gotErr) {
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
			t.Logf("Expected value: %s", wantValueStr)
			t.Logf("Got error:      %s", gotErr)
			t.Fail()
		} else if !types.Equivalent(wantType, gotType) {
			t.Log("ERROR: Program has unexpected output type.")
			t.Logf("Program:        %s", program)
			t.Logf("Expected type:  %s", wantType)
			t.Logf("Expected value: %s", wantValueStr)
			t.Logf("Got type:       %s", gotType)
			t.Logf("Got value:      %s", gotValueStr)
			t.Fail()
		} else if ok, _ := wantValue.Equal(gotValue); !ok {
			t.Log("ERROR: Program has unexpected output value.")
			t.Logf("Program:        %s", program)
			t.Logf("Expected type:  %s", wantType)
			t.Logf("Expected value: %s", wantValueStr)
			t.Logf("Got type:       %s", gotType)
			t.Logf("Got value:      %s", gotValueStr)
			t.Fail()
		}
	}
}
