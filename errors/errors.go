package errors

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/participle/lexer"
)

func Explain(kind string, program string, err error) {
	lexerError, ok := err.(*lexer.Error)
	if !ok {
		panic(err)
	}
	lines := strings.SplitAfter(program, "\n")
	line := lines[lexerError.Pos.Line - 1]
	column := lexerError.Pos.Column
	// TODO shorten long lines
	fmt.Fprintln(os.Stderr, kind, "error at", lexerError.Pos)
	fmt.Fprint(os.Stderr, line)
	if len(line) == 0 || line[len(line) - 1] != '\n' {
		fmt.Fprintln(os.Stderr)
	}
	fmt.Fprint(os.Stderr, strings.Repeat(" ", column - 1))
	fmt.Fprintln(os.Stderr, "^")
	fmt.Fprintln(os.Stderr, lexerError.Message)
}
