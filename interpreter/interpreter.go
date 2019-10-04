package interpreter

import (
	"github.com/texttheater/bach/builtin"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/parser"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

// InterpretString takes a Bach program as a string, interprets it and returns
// the result type and value.
func InterpretString(program string) (types.Type, values.Value, error) {
	// parse
	x, err := parser.Parse(program)
	if err != nil {
		return nil, nil, err
	}
	// type-check
	outputShape, action, err := x.Typecheck(builtin.InitialShape, nil)
	if err != nil {
		return nil, nil, err
	}
	if (types.VoidType{}).Subsumes(outputShape.Type) {
		return nil, nil, errors.E(
			errors.Code(errors.VoidProgram),
			errors.Pos(x.Position()),
		)
	}
	// evaluate
	outputState := action(states.InitialState, nil)
	if outputState.Error != nil {
		return nil, nil, err
	}
	drain(outputShape.Type, outputState.Value)
	return outputShape.Type, outputState.Value, nil
}

func drain(t types.Type, v values.Value) {
	if !types.AnySeqType.Subsumes(t) {
		return
	}
	eType := t.ElementType()
	for e := range v.Iter() {
		drain(eType, e)
	}
}
