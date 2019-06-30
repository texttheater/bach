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
	Name *string
}

func (p TypePattern) Typecheck(inputShape shapes.Shape) (shapes.Shape, types.Type, Matcher, error) {
	// partition the input type and check for impossible match
	intersection, complement := inputShape.Type.Partition(p.Type)
	if (types.VoidType{}).Subsumes(intersection) {
		return shapes.Shape{}, nil, nil, errors.E(
			errors.Code(errors.ImpossibleMatch),
			errors.Pos(p.Pos),
			errors.WantType(inputShape.Type),
			errors.GotType(p.Type),
		)
	}
	// build output shape
	outputShape := shapes.Shape{
		Type:  intersection,
		Stack: inputShape.Stack,
	}
	if p.Name != nil {
		outputShape.Stack = &shapes.FuncerStack{
			Head: shapes.VariableFuncer(p, *p.Name, outputShape.Type),
			Tail: outputShape.Stack,
		}
	}
	// build matcher
	matcher := func(inputState states.State) (*states.VariableStack, bool) {
		// TODO For efficiency, we should check inhabitation of a more
		// general type than p.Type if that is equivalent.
		if inputState.Value.Inhabits(p.Type) {
			varStack := inputState.Stack
			if p.Name != nil {
				varStack = &states.VariableStack{
					Head: states.Variable{
						ID:     p,
						Action: states.SimpleAction(inputState.Value),
					},
					Tail: varStack,
				}
			}
			return varStack, true
		}
		return nil, false
	}
	// return
	return outputShape, complement, matcher, nil
}
