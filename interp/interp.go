package interp

import (
	"github.com/texttheater/bach/builtin"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/grammar"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

// InterpretString takes a Bach program as a string, interprets it and returns
// the result type and value.
func InterpretString(program string) (types.Type, values.Value, error) {
	// parse
	x, err := grammar.Parse(program)
	if err != nil {
		return nil, nil, err
	}
	// type-check
	outputShape, action, err := x.Typecheck(builtin.InitialShape, nil)
	if err != nil {
		return nil, nil, err
	}
	// evaluate
	outputState := action(functions.InitialState, nil) // TODO error handling
	drain(outputShape.Type, outputState.Value)
	return outputShape.Type, outputState.Value, nil
}

func drain(t types.Type, v values.Value) {
	if !(&types.SeqType{types.AnyType}).Subsumes(t) {
		return
	}
	eType := t.ElementType()
	for e := range v.Iter() {
		drain(eType, e)
	}
}
