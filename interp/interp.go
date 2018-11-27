package interp

import (
	//"fmt"
	//"os"

	"github.com/texttheater/bach/builtin"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/grammar"
	"github.com/texttheater/bach/values"
)

// InterpretString takes a Bach program as a string, interprets it and returns
// the result value.
func InterpretString(program string) (values.Value, error) {
	// parse
	x, err := grammar.Parse(program)
	if err != nil {
		return nil, err
	}
	// type-check
	_, action, err := x.Typecheck(builtin.InitialShape, nil)
	if err != nil {
		return nil, err
	}
	// evaluate
	outputState := action(functions.InitialState, nil) // TODO error handling
	//fmt.Fprintf(os.Stderr, "Variables in output state:\n")
	//stack := outputState.Stack
	//for stack != nil {
	//	fmt.Fprintf(os.Stderr, "%s\n", stack.Head.Name)
	//	stack = stack.Tail
	//}
	return outputState.Value, nil
}
