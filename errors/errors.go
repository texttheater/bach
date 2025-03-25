package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type errorAttribute func(err *E)

// SyntaxError builds a syntax error value from a number of error attributes.
// The following functions can be used to create error attributes:
//
//	Code
//	Pos
//	Message
//	WantType
//	GotType
//	GotValue
//	InputType
//	Name
//	ArgNum
//	NumParams
//	ParamNum
//	WantParam
//	GotParam
func SyntaxError(atts ...errorAttribute) error {
	return makeError(SyntaxKind, atts...)
}

// TypeError builds a syntax error value from a number of error attributes.
// The following functions can be used to create error attributes:
//
//	Code
//	Pos
//	Message
//	WantType
//	GotType
//	GotValue
//	InputType
//	Name
//	ArgNum
//	NumParams
//	ParamNum
//	WantParam
//	GotParam
func TypeError(atts ...errorAttribute) error {
	return makeError(TypeKind, atts...)
}

// ValueError builds a syntax error value from a number of error attributes.
// The following functions can be used to create error attributes:
//
//	Code
//	Pos
//	Message
//	WantType
//	GotType
//	GotValue
//	InputType
//	Name
//	ArgNum
//	NumParams
//	ParamNum
//	WantParam
//	GotParam
func ValueError(atts ...errorAttribute) error {
	return makeError(ValueKind, atts...)
}

// UnknownError builds a syntax error value from a number of error attributes.
// The following functions can be used to create error attributes:
//
//	Code
//	Pos
//	Message
//	WantType
//	GotType
//	GotValue
//	InputType
//	Name
//	ArgNum
//	NumParams
//	ParamNum
//	WantParam
//	GotParam
func UnknownError(atts ...errorAttribute) error {
	return makeError(UnknownKind, atts...)
}

func makeError(kind ErrorKind, atts ...errorAttribute) error {
	err := E{
		Kind: &kind,
	}
	e := &err
	for _, att := range atts {
		att(e)
	}
	return e
}

func Code(code ErrorCode) errorAttribute {
	return func(err *E) {
		err.Code = &code
	}
}

func Pos(pos lexer.Position) errorAttribute {
	return func(err *E) {
		err.Pos = &pos
	}
}

func Message(message string) errorAttribute {
	return func(err *E) {
		err.Message = &message
	}
}

func WantType(wantType types.Type) errorAttribute {
	return func(err *E) {
		err.WantType = wantType
	}
}

func GotType(gotType types.Type) errorAttribute {
	return func(err *E) {
		err.GotType = gotType
	}
}

func GotValue(gotValue states.Value) errorAttribute {
	return func(err *E) {
		err.GotValue = gotValue
	}
}

func InputType(inputType types.Type) errorAttribute {
	return func(err *E) {
		err.InputType = inputType
	}
}

func Name(name string) errorAttribute {
	return func(err *E) {
		err.Name = &name
	}
}

func ArgNum(argNum int) errorAttribute {
	return func(err *E) {
		err.ArgNum = &argNum
	}
}

func NumParams(numParams int) errorAttribute {
	return func(err *E) {
		err.NumParams = &numParams
	}
}

func ParamNum(paramNum int) errorAttribute {
	return func(err *E) {
		err.ParamNum = &paramNum
	}
}

func WantParam(wantParam *params.Param) errorAttribute {
	return func(err *E) {
		err.WantParam = wantParam
	}
}

func GotParam(gotParam *params.Param) errorAttribute {
	return func(err *E) {
		err.GotParam = gotParam
	}
}

func Hint(hint string) errorAttribute {
	return func(err *E) {
		err.Hint = &hint
	}
}

// An E represents any code of Bach error, or error template.
type E struct {
	Kind      *ErrorKind
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
	WantParam *params.Param
	GotParam  *params.Param
	Hint      *string
}

func (err *E) Error() string {
	m := make(map[string]any)
	if err.Kind != nil {
		m["Kind"] = err.Kind.String()
	}
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
		m["GotValue"], _ = err.GotValue.Repr()
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
	e, ok := err.(*E)
	if !ok {
		fmt.Fprintln(os.Stderr, "Unknown error")
		fmt.Fprintln(os.Stderr, "Message:   ", err.Error())
		return
	}
	// header and position
	fmt.Fprint(os.Stderr, e.Kind)
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
		gotValueStr, _ := e.GotValue.Repr()
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
	e1, ok := err1.(*E)
	if !ok {
		return false
	}
	e2, ok := err2.(*E)
	if !ok {
		return false
	}
	if e1.Kind != nil && *e2.Kind != *e1.Kind {
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
