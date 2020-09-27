package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type errorAttribute func(err *e)

// E builds an error value from a number of error attributes. The following
// functions can be used to create error attributes:
//
//    Code
//    Pos
//    Message
//    WantType
//    GotType
//    GotValue
//    InputType
//    Name
//    ArgNum
//    NumParams
//    ParamNum
//    WantParam
//    GotParam
func E(atts ...errorAttribute) error {
	err := e{}
	e := &err
	for _, att := range atts {
		att(e)
	}
	return e
}

func Code(code ErrorCode) errorAttribute {
	return func(err *e) {
		err.Code = &code
	}
}

func Pos(pos lexer.Position) errorAttribute {
	return func(err *e) {
		err.Pos = &pos
	}
}

func Message(message string) errorAttribute {
	return func(err *e) {
		err.Message = &message
	}
}

func WantType(wantType types.Type) errorAttribute {
	return func(err *e) {
		err.WantType = wantType
	}
}

func GotType(gotType types.Type) errorAttribute {
	return func(err *e) {
		err.GotType = gotType
	}
}

func GotValue(gotValue states.Value) errorAttribute {
	return func(err *e) {
		err.GotValue = gotValue
	}
}

func InputType(inputType types.Type) errorAttribute {
	return func(err *e) {
		err.InputType = inputType
	}
}

func Name(name string) errorAttribute {
	return func(err *e) {
		err.Name = &name
	}
}

func ArgNum(argNum int) errorAttribute {
	return func(err *e) {
		err.ArgNum = &argNum
	}
}

func NumParams(numParams int) errorAttribute {
	return func(err *e) {
		err.NumParams = &numParams
	}
}

func ParamNum(paramNum int) errorAttribute {
	return func(err *e) {
		err.ParamNum = &paramNum
	}
}

func WantParam(wantParam *parameters.Parameter) errorAttribute {
	return func(err *e) {
		err.WantParam = wantParam
	}
}

func GotParam(gotParam *parameters.Parameter) errorAttribute {
	return func(err *e) {
		err.GotParam = gotParam
	}
}

func Hint(hint string) errorAttribute {
	return func(err *e) {
		err.Hint = &hint
	}
}

// An e represents any code of Bach error, or error template.
type e struct {
	Code      *ErrorCode
	Pos       *lexer.Position
	Message   *string
	WantType  types.Type
	GotType   types.Type
	GotValue  states.Value
	InputType types.Type
	Name      *string
	ArgNum    *int
	NumParams *int
	ParamNum  *int
	WantParam *parameters.Parameter
	GotParam  *parameters.Parameter
	Hint      *string
}

func (err *e) Error() string {
	m := make(map[string]interface{})
	if err.Code != nil {
		m["Code"] = err.Code.String()
	}
	if err.Pos != nil {
		m["Pos"] = err.Pos.String()
	}
	if err.Message != nil {
		m["Message"] = *err.Message
	}
	if err.WantType != nil {
		m["WantType"] = err.WantType.String()
	}
	if err.GotType != nil {
		m["GotType"] = err.GotType.String()
	}
	if err.GotValue != nil {
		m["GotValue"], _ = err.GotValue.String()
	}
	if err.InputType != nil {
		m["InputType"] = err.InputType.String()
	}
	if err.Name != nil {
		m["Name"] = *err.Name
	}
	if err.ArgNum != nil {
		m["ArgNum"] = *err.ArgNum
	}
	if err.NumParams != nil {
		m["NumParams"] = *err.NumParams
	}
	if err.ParamNum != nil {
		m["ParamNum"] = *err.ParamNum
	}
	if err.WantParam != nil {
		m["WantParam"] = err.WantParam.String()
	}
	if err.GotParam != nil {
		m["GotParam"] = err.GotParam.String()
	}
	buffer := bytes.Buffer{}
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	err2 := encoder.Encode(m)
	if err2 != nil {
		panic(err2)
	}
	return buffer.String()
}

func Explain(err error, program string) {
	e, ok := err.(*e)
	if !ok {
		fmt.Fprintln(os.Stderr, "Unknown error")
		fmt.Fprintln(os.Stderr, "Message:   ", err.Error())
		return
	}
	// header and position
	fmt.Fprint(os.Stderr, e.Code.Kind())
	if e.Pos != nil && e.Pos.Line > 0 {
		fmt.Fprintln(os.Stderr, " at", e.Pos)
		lines := strings.SplitAfter(program, "\n")
		line := lines[e.Pos.Line-1]
		column := e.Pos.Column
		// TODO shorten long lines
		fmt.Fprint(os.Stderr, line)
		if len(line) == 0 || line[len(line)-1] != '\n' {
			fmt.Fprintln(os.Stderr)
		}
		fmt.Fprint(os.Stderr, strings.Repeat(" ", column-1))
		fmt.Fprintln(os.Stderr, "^")
	} else {
		fmt.Fprintln(os.Stderr, "")
	}
	// attributes
	fmt.Fprintln(os.Stderr, "Code:      ", e.Code.String())
	if e.Message == nil {
		fmt.Fprintln(os.Stderr, "Message:   ", e.Code.DefaultMessage())
	} else {
		fmt.Fprintln(os.Stderr, "Message:   ", *e.Message)
	}
	if e.WantType != nil {
		fmt.Fprintln(os.Stderr, "Want type: ", e.WantType)
	}
	if e.GotType != nil {
		fmt.Fprintln(os.Stderr, "Got type:  ", e.GotType)
	}
	if e.GotValue != nil {
		gotValueStr, _ := e.GotValue.String()
		fmt.Fprintln(os.Stderr, "Got value: ", gotValueStr)
	}
	if e.InputType != nil {
		fmt.Fprintln(os.Stderr, "Input type:", e.InputType)
	}
	if e.Name != nil {
		fmt.Fprintln(os.Stderr, "Name:      ", *e.Name)
	}
	if e.ArgNum != nil {
		fmt.Fprintln(os.Stderr, "Arg #:     ", *e.ArgNum)
	}
	if e.NumParams != nil {
		fmt.Fprintln(os.Stderr, "# params:  ", *e.NumParams)
	}
	if e.ParamNum != nil {
		fmt.Fprintln(os.Stderr, "Param #:   ", *e.ParamNum)
	}
	if e.WantParam != nil {
		fmt.Fprintln(os.Stderr, "Want param:", e.WantParam)
	}
	if e.GotParam != nil {
		fmt.Fprintln(os.Stderr, "Got param: ", e.GotParam)
	}
}

// Match compares its two error arguments. It can be used to check for expected
// errors in tests. The arguments must both have underlying type *e or
// Match will return false. Otherwise it returns true iff every non-none
// element of the first error is equal to the corresponding element of the
// second. Elements that are in the second argument but not present in the
// first are ignored.
//
// Adapted from: https://github.com/upspin/upspin/blob/master/errors/errors.go
func Match(err1, err2 error) bool {
	e1, ok := err1.(*e)
	if !ok {
		return false
	}
	e2, ok := err2.(*e)
	if !ok {
		return false
	}
	if e1.Code != nil && *e2.Code != *e1.Code {
		return false
	}
	if e1.Pos != nil && *e2.Pos != *e1.Pos {
		return false
	}
	if e1.Message != nil && *e2.Message != *e1.Message {
		return false
	}
	if e1.WantType != nil && !types.Equivalent(e1.WantType, e2.WantType) {
		return false
	}
	if e1.GotType != nil && !types.Equivalent(e1.GotType, e2.GotType) {
		return false
	}
	if e1.GotValue != nil {
		if ok, _ := e1.GotValue.Equal(e2.GotValue); !ok {
			return false
		}
	}
	if e1.InputType != nil && !types.Equivalent(e1.InputType, e2.InputType) {
		return false
	}
	if e1.Name != nil && *e2.Name != *e1.Name {
		return false
	}
	if e1.ArgNum != nil && *e2.ArgNum != *e1.ArgNum {
		return false
	}
	if e1.NumParams != nil && *e2.NumParams != *e1.NumParams {
		return false
	}
	if e1.ParamNum != nil && *e2.ParamNum != *e1.ParamNum {
		return false
	}
	if e1.WantParam != nil && !e1.WantParam.Equivalent(e2.WantParam) {
		return false
	}
	if e1.GotParam != nil && !e1.GotParam.Equivalent(e2.GotParam) {
		return false
	}
	return true
}
