package patterns

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type ObjPattern struct {
	Pos            lexer.Position
	PropPatternMap map[string]Pattern
	Name           *string
}

func (p *ObjPattern) Typecheck(inputShape shapes.Shape) (shapes.Shape, types.Type, Matcher, error) {
	// compute value input types
	propInputTypeMap := make(map[string]types.Type)
	switch t := inputShape.Type.(type) {
	case types.ObjType:
		for prop := range p.PropPatternMap {
			valType, ok := t.PropTypeMap[prop]
			if !ok {
				valType = types.VoidType{}
			}
			propInputTypeMap[prop] = valType
		}
	case types.UnionType:
	PatternProps:
		for prop := range p.PropPatternMap {
			propInputTypeMap[prop] = types.VoidType{}
			for _, disjunct := range t {
				switch d := disjunct.(type) {
				case types.ObjType:
					valType, ok := d.PropTypeMap[prop]
					if !ok {
						propInputTypeMap[prop] = types.AnyType{}
						continue PatternProps
					}
					propInputTypeMap[prop] = types.Union(
						propInputTypeMap[prop],
						valType,
					)
				}
			}
		}
	default:
		for prop := range p.PropPatternMap {
			propInputTypeMap[prop] = types.VoidType{}
		}
	}
	// process value patterns
	funcerStack := inputShape.Stack
	propTypeMap := make(map[string]types.Type)
	propMatcherMap := make(map[string]Matcher)
	for prop, valPattern := range p.PropPatternMap {
		valShape, _, valMatcher, err := valPattern.Typecheck(shapes.Shape{
			Type:  propInputTypeMap[prop],
			Stack: funcerStack,
		})
		if err != nil {
			return shapes.Shape{}, nil, nil, err
		}
		funcerStack = valShape.Stack
		propTypeMap[prop] = valShape.Type
		propMatcherMap[prop] = valMatcher
	}
	// determine the type of values this pattern will match
	pType := types.NewObjType(propTypeMap)
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
		case values.ObjValue:
			for prop, valMatcher := range propMatcherMap {
				value, ok := v[prop]
				if !ok {
					return nil, false
				}
				varStack, ok = valMatcher(states.State{
					Value: value,
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
