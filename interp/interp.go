package interp

import (
	"github.com/texttheater/bach/builtin"
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
	f, err := x.Function(builtin.InitialShape)
	if err != nil {
		return nil, err
	}
	// evaluate
	return f.OutputState(builtin.InitialState).Value, nil // TODO error handling
}