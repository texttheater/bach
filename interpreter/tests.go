package interpreter

import (
	"log"
	"math"
	"testing"

	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/grammar"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

func TestExample(example shapes.Example) {
	program := example.Program
	wantTypeStr := example.OutputType
	wantValueStr := example.OutputValue
	wantError := example.Error
	var wantType types.Type
	var err error
	if example.OutputType != "" {
		wantType, err = grammar.ParseType(wantTypeStr)
		if err != nil {
			log.Println("ERROR: Could not parse expected type")
			errors.Explain(err, wantTypeStr)
			log.Fatal()
		}
	}
	var wantValue states.Value
	if example.OutputValue != "" {
		_, wantValue, err = InterpretString(wantValueStr)
		if err != nil {
			log.Println("ERROR: Could not interpret expected value")
			errors.Explain(err, wantValueStr)
			log.Fatal()
		}
	}
	gotType, gotValue, gotErr := InterpretString(example.Program)
	var gotValueStr string
	if gotValue != nil {
		var err error
		gotValueStr, err = gotValue.Repr()
		if gotErr == nil {
			gotErr = err
		}
	}
	if wantError != nil {
		if gotErr == nil {
			log.Print("ERROR: Expected error but program succeeded.")
			log.Printf("Program:        %s", program)
			log.Printf("Expected error: %s", wantError)
			log.Printf("Got type:       %s", gotType)
			log.Printf("Got value:      %s", gotValueStr)
			log.Fatal()
		} else if !errors.Match(wantError, gotErr) {
			log.Printf("ERROR: Expected error does not match actual error.")
			log.Printf("Program:        %s", program)
			log.Printf("Expected error: %s", wantError)
			log.Printf("Got error:      %s", gotErr)
			log.Fatal()
		}
	} else {
		if gotErr != nil {
			log.Print("ERROR: Expected program to succeed but got error.")
			log.Printf("Program:        %s", program)
			log.Printf("Expected type:  %s", wantType)
			log.Printf("Expected value: %s", wantValueStr)
			log.Printf("Got error:      %s", gotErr)
			log.Fatal()
		} else if !types.Equivalent(wantType, gotType) {
			log.Print("ERROR: Program has unexpected output type.")
			log.Printf("Program:        %s", program)
			log.Printf("Expected type:  %s", wantType)
			log.Printf("Expected value: %s", wantValueStr)
			log.Printf("Got type:       %s", gotType)
			log.Printf("Got value:      %s", gotValueStr)
			log.Fatal()
		} else if ok, _ := match(wantValue, gotValue); !ok {
			log.Print("ERROR: Program has unexpected output value.")
			log.Printf("Program:        %s", program)
			log.Printf("Expected type:  %s", wantType)
			log.Printf("Expected value: %s", wantValueStr)
			log.Printf("Got type:       %s", gotType)
			log.Printf("Got value:      %s", gotValueStr)
			log.Fatal()
		}
	}
}

func TestProgramStr(program string, wantTypeString string, wantValueString string, wantError error, t *testing.T) {
	var wantType types.Type
	var err error
	if wantTypeString != "" {
		wantType, err = grammar.ParseType(wantTypeString)
		if err != nil {
			t.Log("ERROR: Could not parse expected type")
			errors.Explain(err, wantTypeString)
			t.Fail()
		}
	}
	var wantValue states.Value
	if wantValueString != "" {
		_, wantValue, err = InterpretString(wantValueString)
		if err != nil {
			t.Log("ERROR: Could not interpret expected value")
			errors.Explain(err, wantValueString)
			t.Fail()
		}
	}
	TestProgram(program, wantType, wantValue, wantError, t)
}

func TestProgram(program string, wantType types.Type, wantValue states.Value, wantError error, t *testing.T) {
	//log.Print(program)
	var wantValueStr string
	if wantValue != nil {
		wantValueStr, _ = wantValue.Repr()
	}
	gotType, gotValue, gotErr := InterpretString(program)
	var gotValueStr string
	if gotValue != nil {
		var err error
		gotValueStr, err = gotValue.Repr()
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
		} else if ok, _ := match(wantValue, gotValue); !ok {
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

func match(want states.Value, got states.Value) (bool, error) {
	w, ok := want.(states.NumValue)
	if ok {
		g, ok := got.(states.NumValue)
		if ok {
			if math.IsNaN(float64(w)) && math.IsNaN(float64(g)) {
				return true, nil
			}
		}
	}
	return want.Equal(got)
}
