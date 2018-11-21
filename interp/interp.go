package interp

import (
	//"fmt"
	//"os"

	"github.com/texttheater/bach/builtin"
	"github.com/texttheater/bach/grammar"
	"github.com/texttheater/bach/values"
)

var debug bool = true

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
	outputState := action.Execute(builtin.InitialState, nil)
	return outputState.Value, nil // TODO error handling
}
