package patterns

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type ArrPattern struct {
	Pos             lexer.Position
	ElementPatterns []Pattern
	Name            *string
}

func (p *ArrPattern) Typecheck(inputShape shapes.Shape) (shapes.Shape, types.Type, Matcher, error) {
	// compute element input types
	elementInputTypes := make([]types.Type, len(p.ElementPatterns))
	switch t := inputShape.Type.(type) {
	case types.TupType:
		if len(t) != len(elementInputTypes) {
			return shapes.Shape{}, types.VoidType{}, nil, nil
		}
		for i, elType := range t {
			elementInputTypes[i] = elType
		}
	case *types.ArrType:
		for i := range elementInputTypes {
			elementInputTypes[i] = t.ElType
		}
	case types.UnionType:
		for _, disjunct := range t {
			switch d := disjunct.(type) {
			case types.TupType:
				if len(d) != len(elementInputTypes) {
					continue
				}
				for i := range elementInputTypes {
					elementInputTypes[i] = types.Union(
						elementInputTypes[i],
						d[i],
					)
				}
			case *types.ArrType:
				for i := range elementInputTypes {
					elementInputTypes[i] = types.Union(
						elementInputTypes[i],
						d.ElType,
					)
				}
			}
		}
		for _, elInputType := range elementInputTypes {
			if (types.VoidType{}).Subsumes(elInputType) {
				return shapes.Shape{}, types.VoidType{}, nil, nil
			}
		}
	default:
		return shapes.Shape{}, types.VoidType{}, nil, nil
	}
	// process element patterns
	funcerStack := inputShape.Stack
	elementTypes := make([]types.Type, len(p.ElementPatterns))
	elementMatchers := make([]Matcher, len(p.ElementPatterns))
	for i, elPattern := range p.ElementPatterns {
		elShape, _, elMatcher, err := elPattern.Typecheck(shapes.Shape{
			Type:  elementInputTypes[i],
			Stack: funcerStack,
		})
		if err != nil {
			return shapes.Shape{}, nil, nil, err
		}
		funcerStack = elShape.Stack
		elementTypes[i] = elShape.Type
		elementMatchers[i] = elMatcher
	}
	// determine the type of values this pattern will match
	pType := types.TupType(elementTypes)
	// partition the input type and check for impossible match
	intersection, complement := inputShape.Type.Partition(pType)
	if (types.VoidType{}).Subsumes(intersection) {
		return shapes.Shape{}, nil, nil, errors.E(
			errors.Code(errors.ImpossibleMatch),
			errors.Pos(p.Pos),
			errors.WantType(inputShape.Type),
			errors.GotType(pType),
		)
	}
	// build output shape
	outputShape := shapes.Shape{
		Type:  intersection,
		Stack: funcerStack,
	}
	if p.Name != nil {
		outputShape.Stack = &shapes.FuncerStack{
			Head: shapes.VariableFuncer(p, *p.Name, outputShape.Type),
			Tail: outputShape.Stack,
		}
	}
	// build matcher
	matcher := func(inputState states.State) (*states.VariableStack, bool) {
		varStack := inputState.Stack
		switch v := inputState.Value.(type) {
		case values.ArrValue:
			if len(v) != len(elementMatchers) {
				return nil, false
			}
			for i, elMatcher := range elementMatchers {
				var ok bool
				varStack, ok = elMatcher(states.State{
					Value: v[i],
					Stack: varStack,
				})
				if !ok {
					return nil, false
				}
			}
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
		default:
			return nil, false
		}
	}
	// return
	return outputShape, complement, matcher, nil

}
