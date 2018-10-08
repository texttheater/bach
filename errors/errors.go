package errors

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/participle/lexer"
)

type Error struct {
	Kind    string
	Pos     lexer.Position
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func E(kind string, pos lexer.Position, format string, a ...interface{}) error {
	return Error{kind, pos, fmt.Sprintf(format, a...)}
}

func Is(kind string, err error) bool {
	e, ok := err.(Error)
	if !ok {
		return false
	}
	return e.Kind == kind
}

func Explain(err error, program string) {
	e, ok := err.(Error)
	if !ok {
		panic(err)
	}
	if e.Pos.Line > 0 && e.Pos.Column > 0 {
		lines := strings.SplitAfter(program, "\n")
		line := lines[e.Pos.Line-1]
		column := e.Pos.Column
		// TODO shorten long lines
		fmt.Fprintln(os.Stderr, e.Kind, "error at", e.Pos)
		fmt.Fprint(os.Stderr, line)
		if len(line) == 0 || line[len(line)-1] != '\n' {
			fmt.Fprintln(os.Stderr)
		}
		fmt.Fprint(os.Stderr, strings.Repeat(" ", column-1))
		fmt.Fprintln(os.Stderr, "^")
	}
	fmt.Fprintln(os.Stderr, e.Message)
}
