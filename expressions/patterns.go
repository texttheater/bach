package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

// pattern/matcher kinda analogous to expression/action

type Pattern interface {
	Position() lexer.Position
	Typecheck(inputShape Shape) (outputShape Shape, restType types.Type, matcher Matcher, err error)
}

type Matcher func(states.State) (*states.VariableStack, bool, error)

type ArrPattern struct {
	Pos             lexer.Position
	ElementPatterns []Pattern
	RestPattern     Pattern
}

func (p ArrPattern) Position() lexer.Position {
	return p.Pos
}

// spreadInputType spreads the input type for an array pattern over its
// elements and rest.
func spreadInputType(inputType types.Type, elementTypes []types.Type) (restType types.Type, ok bool) {
	switch t := inputType.(type) {
	case *types.Nearr:
		if len(elementTypes) == 0 {
			return t, true
		}
		elementTypes[0] = t.Head
		return spreadInputType(t.Tail, elementTypes[1:])
	case *types.Arr:
		// Optional: fail if the pattern wants to match more elements
		// than the value can contain, as per its type. For now, it is
		// is commented out and will instead lead to an error message
		// about a surplus element not having type Void. That's a bit
		// opaque but has the advantage of indicating the place where
		// the array pattern is too long.
		//if (types.VoidType{}).Subsumes(t.ElType) && len(elementTypes) > 0 {
		//	return nil, false
		//}
		for i := range elementTypes {
			elementTypes[i] = t.El
		}
		return t, true
	case types.Union:
		for i := range elementTypes {
			elementTypes[i] = types.Void{}
		}
		var restType types.Type = types.Void{}
		anyOk := false
		for _, disjunct := range t {
			disjunctElementTypes := make([]types.Type, len(elementTypes))
			disjunctRestType, ok := spreadInputType(disjunct, disjunctElementTypes)
			if ok {
				for i := range elementTypes {
					elementTypes[i] = types.NewUnion(elementTypes[i], disjunctElementTypes[i])
				}
				restType = types.NewUnion(restType, disjunctRestType)
				anyOk = true
			}
		}
		return restType, anyOk
	default:
		return nil, false
	}
}

func (p ArrPattern) Typecheck(inputShape Shape) (Shape, types.Type, Matcher, error) {
	// compute element and rest input types
	elementInputTypes := make([]types.Type, len(p.ElementPatterns))
	restInputType, ok := spreadInputType(inputShape.Type, elementInputTypes)
	if !ok {
		return Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.ImpossibleMatch),
			errors.Pos(p.Pos),
		)
	}
	// process element patterns
	funcerStack := inputShape.Stack
	elementTypes := make([]types.Type, len(p.ElementPatterns))
	elementMatchers := make([]Matcher, len(p.ElementPatterns))
	for i, elPattern := range p.ElementPatterns {
		elShape, _, elMatcher, err := elPattern.Typecheck(Shape{
			Type:  elementInputTypes[i],
			Stack: funcerStack,
		})
		if err != nil {
			return Shape{}, nil, nil, err
		}
		funcerStack = elShape.Stack
		elementTypes[i] = elShape.Type
		elementMatchers[i] = elMatcher
	}
	// process rest pattern
	restShape, _, restMatcher, err := p.RestPattern.Typecheck(Shape{
		Type:  restInputType,
		Stack: funcerStack,
	})
	if err != nil {
		return Shape{}, nil, nil, err
	}
	funcerStack = restShape.Stack
	// determine the type of values this pattern will match
	pType := types.NewNearr(elementTypes, restShape.Type)
	// partition the input type and check for impossible match
	intersection, complement := inputShape.Type.Partition(pType)
	if (types.Void{}).Subsumes(intersection) {
		return Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.ImpossibleMatch),
			errors.Pos(p.Pos),
			errors.WantType(inputShape.Type),
			errors.GotType(pType),
		)
	}
	// build output shape
	outputShape := Shape{
		Type:  intersection,
		Stack: funcerStack,
	}
	// build matcher
	matcher := func(inputState states.State) (*states.VariableStack, bool, error) {
		varStack := inputState.Stack
		switch v := inputState.Value.(type) {
		case *states.ArrValue:
			for _, elMatcher := range elementMatchers {
				if v == nil {
					return nil, false, nil
				}
				var ok bool
				varStack, ok, err = elMatcher(states.State{
					Value:     v.Head,
					Stack:     varStack,
					TypeStack: inputState.TypeStack,
				})
				if !ok {
					return nil, false, err
				}
				v, err = v.Tail.EvalArr()
				if err != nil {
					return nil, false, err
				}
			}
			varStack, ok, err = restMatcher(states.State{
				Value:     v,
				Stack:     varStack,
				TypeStack: inputState.TypeStack,
			})
			if !ok {
				return nil, false, err
			}
			return varStack, true, nil
		default:
			return nil, false, nil
		}
	}
	// return
	return outputShape, complement, matcher, nil
}

type ObjPattern struct {
	Pos            lexer.Position
	PropPatternMap map[string]Pattern
}

func (p ObjPattern) Position() lexer.Position {
	return p.Pos
}

func (p ObjPattern) Typecheck(inputShape Shape) (Shape, types.Type, Matcher, error) {
	// compute value input types
	propInputTypeMap := make(map[string]types.Type)
	switch t := inputShape.Type.(type) {
	case types.Obj:
		for prop := range p.PropPatternMap {
			valType, ok := t.Props[prop]
			if !ok {
				valType = t.Rest
			}
			propInputTypeMap[prop] = valType
		}
	case types.Union:
	PatternProps:
		for prop := range p.PropPatternMap {
			propInputTypeMap[prop] = types.Void{}
			for _, disjunct := range t {
				switch d := disjunct.(type) {
				case types.Obj:
					valType, ok := d.Props[prop]
					if !ok {
						propInputTypeMap[prop] = types.Any{}
						continue PatternProps
					}
					propInputTypeMap[prop] = types.NewUnion(
						propInputTypeMap[prop],
						valType,
					)
				}
			}
		}
	case types.Any:
		for prop := range p.PropPatternMap {
			propInputTypeMap[prop] = types.Any{}
		}
	default:
		for prop := range p.PropPatternMap {
			propInputTypeMap[prop] = types.Void{}
		}
	}
	// process value patterns
	funcerStack := inputShape.Stack
	propTypeMap := make(map[string]types.Type)
	propMatcherMap := make(map[string]Matcher)
	for prop, valPattern := range p.PropPatternMap {
		valShape, _, valMatcher, err := valPattern.Typecheck(Shape{
			Type:  propInputTypeMap[prop],
			Stack: funcerStack,
		})
		if err != nil {
			return Shape{}, nil, nil, err
		}
		funcerStack = valShape.Stack
		propTypeMap[prop] = valShape.Type
		propMatcherMap[prop] = valMatcher
	}
	// determine the type of values this pattern will match
	pType := types.Obj{
		Props: propTypeMap,
		Rest:  types.Any{},
	}
	// partition the input type and check for impossible match
	intersection, complement := inputShape.Type.Partition(pType)
	if (types.Void{}).Subsumes(intersection) {
		return Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.ImpossibleMatch),
			errors.Pos(p.Pos),
			errors.WantType(inputShape.Type),
			errors.GotType(pType),
		)
	}
	// build output shape
	outputShape := Shape{
		Type:  intersection,
		Stack: funcerStack,
	}
	// build matcher
	matcher := func(inputState states.State) (*states.VariableStack, bool, error) {
		varStack := inputState.Stack
		switch v := inputState.Value.(type) {
		case states.ObjValue:
			for prop, valMatcher := range propMatcherMap {
				thunk, ok := v[prop]
				if !ok {
					return nil, false, nil
				}
				val, err := thunk.Eval()
				if err != nil {
					return nil, false, err
				}
				varStack, ok, err = valMatcher(states.State{
					Value:     val,
					Stack:     varStack,
					TypeStack: inputState.TypeStack,
				})
				if !ok {
					return nil, false, err
				}
			}
			return varStack, true, nil
		default:
			return nil, false, nil
		}
	}
	// return
	return outputShape, complement, matcher, nil
}

type TypePattern struct {
	Pos  lexer.Position
	Type types.Type
	Name *string
}

func (p TypePattern) Position() lexer.Position {
	return p.Pos
}

func (p TypePattern) Typecheck(inputShape Shape) (Shape, types.Type, Matcher, error) {
	// partition the input type and check for impossible match
	intersection, complement := inputShape.Type.Partition(p.Type)
	if (types.Void{}).Subsumes(intersection) {
		return Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.ImpossibleMatch),
			errors.Pos(p.Pos),
			errors.WantType(inputShape.Type),
			errors.GotType(p.Type),
		)
	}
	// build output shape
	outputShape := Shape{
		Type:  intersection,
		Stack: inputShape.Stack,
	}
	if p.Name != nil {
		outputShape.Stack = outputShape.Stack.Push(
			VariableFuncer(p, *p.Name, outputShape.Type),
		)
	}
	// build matcher
	matcher := func(inputState states.State) (*states.VariableStack, bool, error) {
		// TODO For efficiency, we should check inhabitation of a more
		// general type than p.Type if that is equivalent.
		if ok, err := inputState.Value.Inhabits(p.Type, inputState.TypeStack); !ok {
			return nil, false, err
		}
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
		return varStack, true, nil
	}
	// return
	return outputShape, complement, matcher, nil
}
