package patterns

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type TypePattern struct {
	Pos  lexer.Position
	Type types.Type
}

func (p TypePattern) Typecheck(inputShape shapes.Shape) (shapes.Shape, types.Type, Matcher, error) {
	intersection, complement := inputShape.Type.Partition(p.Type)
	if (types.VoidType{}).Subsumes(intersection) {
		return shapes.Shape{}, nil, nil, errors.E(
			errors.Code(errors.ImpossibleMatch),
			errors.Pos(p.Pos),
			errors.WantType(inputShape.Type),
			errors.GotType(p.Type),
		)
	}
	outputShape := shapes.Shape{
		Type:  intersection,
		Stack: inputShape.Stack,
	}
	matcher := func(inputState states.State) (*states.VariableStack, bool) {
		// TODO For efficiency, we should check inhabitation of a more
		// general type than p.Type if that is equivalent.
		if inputState.Value.Inhabits(p.Type) {
			return inputState.Stack, true
		}
		return nil, false
	}
	return outputShape, complement, matcher, nil
}
