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
	_, action, err := x.Typecheck(builtin.InitialContext, nil)
	if err != nil {
		return nil, err
	}
	// evaluate
	return action(&values.NullValue{}, nil), nil // TODO error handling
}
