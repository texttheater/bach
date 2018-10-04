package errors

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/participle/lexer"
)

type bachError struct {
	cat string
	pos lexer.Position
	msg string
}

func (be bachError) Error() string {
	return be.msg
}

func Errorf(cat string, pos lexer.Position, format string, a ...interface{}) bachError {
	return bachError{cat, pos, fmt.Sprintf(format, a)}
}

func Explain(err error, program string) {
	be, ok := err.(bachError)
	if ok {
		explain(be.cat, be.pos, be.msg, program)
		return
	}
	lexerError, ok := err.(*lexer.Error)
	if ok {
		explain("syntax", lexerError.Pos, lexerError.Message, program)
		return
	}
	panic(err)
}

func explain(cat string, pos lexer.Position, msg string, program string) {
	if pos.Line > 0 && pos.Column > 0 {
		lines := strings.SplitAfter(program, "\n")
		line := lines[pos.Line-1]
		column := pos.Column
		// TODO shorten long lines
		fmt.Fprintln(os.Stderr, cat, "error at", pos)
		fmt.Fprint(os.Stderr, line)
		if len(line) == 0 || line[len(line)-1] != '\n' {
			fmt.Fprintln(os.Stderr)
		}
		fmt.Fprint(os.Stderr, strings.Repeat(" ", column-1))
		fmt.Fprintln(os.Stderr, "^")
	}
	fmt.Fprintln(os.Stderr, msg)
}
