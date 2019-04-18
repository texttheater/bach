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
	if !inputShape.Type.Subsumes(p.Type) {
		return zeroShape, nil, nil, errors.E(
			errors.Code(errors.ImpossibleMatch),
			errors.Pos(p.Pos),
			errors.WantType(inputShape.Type),
			errors.GotType(p.Type),
		)
	}
	outputShape := shapes.Shape{
		Type:        p.Type,
		FuncerStack: inputShape.FuncerStack,
	}
	restType := types.Complement(inputShape.Type, p.Type)
	matcher := func(inputState states.State) (*states.VariableStack, bool) {
		if inputState.Value.Inhabits(p.Type) {
			return inputState.Stack, true
		}
		return nil, false
	}
	return outputShape, restType, matcher, nil
}
