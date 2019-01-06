package errors

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

///////////////////////////////////////////////////////////////////////////////

type ErrorKind int

const (
	Syntax ErrorKind = iota
	ParamsNotAllowed
	NoSuchFunction
	ArgHasWrongOutputType
	ParamDoesNotMatch
	FunctionBodyHasWrongOutputType
	ConditionMustBeBool
	MappingRequiresSeqType
)

func (kind ErrorKind) String() string {
	switch kind {
	case Syntax:
		return "Syntax"
	case ParamsNotAllowed:
		return "ParamsNotAllowed"
	case NoSuchFunction:
		return "NoSuchFunction"
	case ArgHasWrongOutputType:
		return "ArgHasWrongOutputType"
	case FunctionBodyHasWrongOutputType:
		return "FunctionBodyHasWrongOutputType"
	case ConditionMustBeBool:
		return "ConditionMustBeBool"
	case MappingRequiresSeqType:
		return "MappingRequiresSeqType"
	}
	return "Unknown"
}

func (kind ErrorKind) DefaultMessage() string {
	switch kind {
	case Syntax:
		return "syntax error"
	case ParamsNotAllowed:
		return "This expression cannot be used as an argument here because it does not take parameters."
	case NoSuchFunction:
		return "no such function"
	case ArgHasWrongOutputType:
		return "An argument has the wrong output type."
	case ParamDoesNotMatch:
		return "Cannot use this function here because one of its parameters does not match the expected parameter."
	case FunctionBodyHasWrongOutputType:
		return "The function body has the wrong output type."
	case ConditionMustBeBool:
		return "The condition must be boolean."
	case MappingRequiresSeqType:
		return "The input to a mapping must be a sequence."
	}
	return "unknown error"
}

///////////////////////////////////////////////////////////////////////////////

type errorAttribute func(err *e)

// E builds an error value from a number of error attributes. The following
// functions can be used to create error attributes:
//
//    Kind
//    Pos
//    Message
//    WantType
//    GotType
//    InputType
//    Name
//    ArgNum
//    NumParams
//    ParamNum
//    WantParam
//    GotParam
func E(atts ...errorAttribute) error {
	err := no
	e := &err
	for _, att := range atts {
		att(e)
	}
	return e
}

func Kind(kind ErrorKind) errorAttribute {
	return func(err *e) {
		err.Kind = kind
	}
}

func Pos(pos lexer.Position) errorAttribute {
	return func(err *e) {
		err.Pos = pos
	}
}

func Message(message string) errorAttribute {
	return func(err *e) {
		err.Message = message
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

func InputType(inputType types.Type) errorAttribute {
	return func(err *e) {
		err.InputType = inputType
	}
}

func Name(name string) errorAttribute {
	return func(err *e) {
		err.Name = name
	}
}

func ArgNum(argNum int) errorAttribute {
	return func(err *e) {
		err.ArgNum = argNum
	}
}

func NumParams(numParams int) errorAttribute {
	return func(err *e) {
		err.NumParams = numParams
	}
}

func ParamNum(numParams int) errorAttribute {
	return func(err *e) {
		err.ParamNum = numParams
	}
}

func WantParam(wantParam *functions.Parameter) errorAttribute {
	return func(err *e) {
		err.WantParam = wantParam
	}
}

func GotParam(gotParam *functions.Parameter) errorAttribute {
	return func(err *e) {
		err.GotParam = gotParam
	}
}

///////////////////////////////////////////////////////////////////////////////

// An e represents any kind of Bach error, or error template. Every field
// may have a "none" value, which is Go's zero value except for int fields
// where it's -1.
type e struct {
	Kind      ErrorKind
	Pos       lexer.Position
	Message   string
	WantType  types.Type
	GotType   types.Type
	InputType types.Type
	Name      string
	ArgNum    int
	NumParams int
	ParamNum  int
	WantParam *functions.Parameter
	GotParam  *functions.Parameter
}

func (err *e) Error() string {
	m := make(map[string]interface{})
	if err.Kind != no.Kind {
		m["Kind"] = err.Kind.String()
	}
	if err.Pos != no.Pos {
		m["Pos"] = err.Pos.String()
	}
	if err.Message != no.Message {
		m["Message"] = err.Message
	}
	if err.WantType != no.WantType {
		m["WantType"] = err.WantType.String()
	}
	if err.GotType != no.GotType {
		m["GotType"] = err.GotType.String()
	}
	if err.InputType != no.InputType {
		m["InputType"] = err.InputType.String()
	}
	if err.Name != no.Name {
		m["Name"] = err.Name
	}
	if err.ArgNum != no.ArgNum {
		m["ArgNum"] = err.ArgNum
	}
	if err.NumParams != no.NumParams {
		m["NumParams"] = err.NumParams
	}
	if err.ParamNum != no.ParamNum {
		m["ParamNum"] = err.ParamNum
	}
	if err.WantParam != no.WantParam {
		m["WantParam"] = err.WantParam.String()
	}
	if err.GotParam != no.GotParam {
		m["GotParam"] = err.GotParam.String()
	}
	out, err2 := json.Marshal(m)
	if err2 != nil {
		panic(err2)
	}
	return string(out)
}

///////////////////////////////////////////////////////////////////////////////

func Explain(err error, program string) {
	e, ok := err.(*e)
	if !ok {
		fmt.Fprintln(os.Stderr, "ERROR")
		fmt.Fprintln(os.Stderr, "Message:   ", err.Error())
		return
	}
	// header and position
	fmt.Fprint(os.Stderr, "ERROR")
	if e.Pos != no.Pos {
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
	if e.Message == no.Message {
		fmt.Fprintln(os.Stderr, "Message:   ", e.Kind.DefaultMessage())
	} else {
		fmt.Fprintln(os.Stderr, "Message:   ", e.Message)
	}
	if e.WantType != no.WantType {
		fmt.Fprintln(os.Stderr, "Want type: ", e.WantType)
	}
	if e.GotType != no.GotType {
		fmt.Fprintln(os.Stderr, "Got type:  ", e.GotType)
	}
	if e.InputType != no.InputType {
		fmt.Fprintln(os.Stderr, "Input type:", e.InputType)
	}
	if e.Name != no.Name {
		fmt.Fprintln(os.Stderr, "Name:      ", e.Name)
	}
	if e.ArgNum != no.ArgNum {
		fmt.Fprintln(os.Stderr, "Arg #:     ", e.ArgNum)
	}
	if e.NumParams != no.NumParams {
		fmt.Fprintln(os.Stderr, "# params:  ", e.NumParams)
	}
	if e.ParamNum != no.ParamNum {
		fmt.Fprintln(os.Stderr, "Param #:   ", e.ParamNum)
	}
	if e.WantParam != no.WantParam {
		fmt.Fprintln(os.Stderr, "Want param:", e.WantParam)
	}
	if e.GotParam != no.GotParam {
		fmt.Fprintln(os.Stderr, "Got param: ", e.GotParam)
	}
}

// no is an e where every field has its "none" value, for convenience.
var no = e{
	Kind:      -1,
	ArgNum:    -1,
	NumParams: -1,
	ParamNum:  -1,
}

///////////////////////////////////////////////////////////////////////////////

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
	if e1.Kind != no.Kind && e2.Kind != e1.Kind {
		return false
	}
	if e1.Pos != no.Pos && e2.Pos != e1.Pos {
		return false
	}
	if e1.Message != no.Message && e2.Message != e1.Message {
		return false
	}
	if e1.WantType != no.WantType && !reflect.DeepEqual(e1.WantType, e1.WantType) {
		return false
	}
	if e1.GotType != no.GotType && !reflect.DeepEqual(e1.GotType, e2.GotType) {
		return false
	}
	if e1.InputType != no.InputType && !reflect.DeepEqual(e1.InputType, e2.InputType) {
		return false
	}
	if e1.Name != no.Name && e2.Name != e1.Name {
		return false
	}
	if e1.ArgNum != no.ArgNum && e2.ArgNum != e1.ArgNum {
		return false
	}
	if e1.NumParams != no.NumParams && e2.NumParams != e1.NumParams {
		return false
	}
	if e1.ParamNum != no.ParamNum && e2.ParamNum != e1.ParamNum {
		return false
	}
	if e1.WantParam != no.WantParam && !reflect.DeepEqual(e1.WantParam, e2.WantParam) {
		return false
	}
	if e1.GotParam != no.GotParam && !reflect.DeepEqual(e1.GotParam, e2.GotParam) {
		return false
	}
	return true
}

////////////////////////////////////////////////////////////////////////////////
