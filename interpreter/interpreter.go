package interpreter

import (
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/grammar"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

// InterpretString takes a Bach program as a string, interprets it and returns
// the result type and value.
func InterpretString(inputShape shapes.Shape, inputState states.State, program string) (types.Type, states.Value, error) {
	// parse
	x, err := grammar.ParseComposition(program)
	if err != nil {
		return nil, nil, err
	}
	// type-check
	outputShape, action, _, err := x.Typecheck(inputShape, nil)
	if err != nil {
		return nil, nil, err
	}
	if (types.Void{}).Subsumes(outputShape.Type) {
		return nil, nil, errors.TypeError(
			errors.Code(errors.VoidProgram),
			errors.Pos(x.Position()),
		)
	}
	// evaluate
	val, err := action(inputState, nil).Eval()
	if err != nil {
		return nil, nil, err
	}
	return outputShape.Type, val, nil
}
