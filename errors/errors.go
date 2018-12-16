package errors

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
)

type E struct {
	Kind      *Kind
	Pos       *lexer.Position
	Message   *string
	WantType  types.Type
	GotType   types.Type
	InputType types.Type
	Name      *string
	NumArgs   *uint8
	ArgNum    *uint8
	NumParams *uint8
	WantParam *functions.Parameter
	GotParam  *functions.Parameter
}

func (err *E) Explain(program string) {
	// header and position
	fmt.Fprint(os.Stderr, "ERROR")
	if err.Pos.Line > 0 && err.Pos.Column > 0 {
		fmt.Fprintln(os.Stderr, " at", err.Pos)
		lines := strings.SplitAfter(program, "\n")
		line := lines[err.Pos.Line-1]
		column := err.Pos.Column
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
	if err.Message == nil {
		fmt.Fprintln(os.Stderr, "Message:   ", err.Kind.DefaultMessage())
	} else {
		fmt.Fprintln(os.Stderr, "Message:   ", err.Message)
	}
	if err.WantType != nil {
		fmt.Fprintln(os.Stderr, "Want type: ", err.WantType)
	}
	if err.GotType != nil {
		fmt.Fprintln(os.Stderr, "Got type:  ", err.GotType)
	}
	if err.InputType != nil {
		fmt.Fprintln(os.Stderr, "Input type:", err.InputType)
	}
	if err.NumArgs != nil {
		fmt.Fprintln(os.Stderr, "# args:    ", err.NumArgs)
	}
	if err.ArgNum != nil {
		fmt.Fprintln(os.Stderr, "Arg #:     ", err.ArgNum)
	}
	if err.NumParams != nil {
		fmt.Fprintln(os.Stderr, "# params:  ", err.NumArgs)
	}
	if err.WantParam != nil {
		fmt.Fprintln(os.Stderr, "Want param:", err.WantParam)
	}
	if err.GotParam != nil {
		fmt.Fprintln(os.Stderr, "Got param: ", err.GotParam)
	}
}

type Kind uint8

const (
	Syntax Kind = iota
	ParamsNotAllowed
	NoSuchFunction
	ArgHasWrongOutputType
	ArgDoesNotMatchParam
	FunctionBodyHasWrongOutputType
	ConditionMustBeBool
)

func (kind Kind) DefaultMessage() string {
	switch kind {
	case Syntax:
		return "syntax error"
	case ParamsNotAllowed:
		return "This expression cannot be used as an argument here because it does not take parameters."
	case NoSuchFunction:
		return "No such function."
	case ArgHasWrongOutputType:
		return "An argument has the wrong output type."
	case FunctionBodyHasWrongOutputType:
		return "The function body has the wrong output type."
	case ConditionMustBeBool:
		return "The condition must be boolean."
	}
	return "unknown error"
}

func (err *E) Subsumes(other *E) bool {
	if err.Kind != nil && err.Kind != other.Kind {
		return false
	}
	if err.Pos != nil && err.Pos != other.Pos {
		return false
	}
	if err.Message != nil && err.Message != other.Message {
		return false
	}
	if err.WantType != nil && err.WantType != other.WantType {
		return false
	}
	if err.GotType != nil && err.GotType != other.GotType {
		return false
	}
	if err.InputType != nil && err.InputType != other.InputType {
		return false
	}
	if err.Name != nil && err.Name != other.Name {
		return false
	}
	if err.NumArgs != nil && err.NumArgs != other.NumArgs {
		return false
	}
	if err.WantParam != nil && err.WantParam != other.WantParam {
		return false
	}
	if err.GotParam != nil && err.GotParam != other.GotParam {
		return false
	}
	return true
}
