package interpreter

import (
	"github.com/texttheater/bach/builtin"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/grammar"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

// InterpretString takes a Bach program as a string, interprets it and returns
// the result type and value.
func InterpretString(program string) (types.Type, states.Value, error) {
	// parse
	x, err := grammar.ParseComposition(program)
	if err != nil {
		return nil, nil, err
	}
	// type-check
	outputShape, action, _, err := x.Typecheck(builtin.InitialShape, nil)
	if err != nil {
		return nil, nil, err
	}
	if (types.VoidType{}).Subsumes(outputShape.Type) {
		return nil, nil, errors.E(
			errors.Code(errors.VoidProgram),
			errors.Pos(x.Position()))

	}
	// evaluate
	res := action(states.InitialState, nil).Eval()
	if res.Error != nil {
		return nil, nil, res.Error
	}
	return outputShape.Type, res.Value, nil
}
