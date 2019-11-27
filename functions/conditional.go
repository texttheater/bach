package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type ConditionalExpression struct {
	Pos                           lexer.Position
	Pattern                       Pattern
	Guard                         Expression
	Consequent                    Expression
	AlternativePatterns           []Pattern
	AlternativeGuards             []Expression
	AlternativeConsequents        []Expression
	Alternative                   Expression
	UnreachableAlternativeAllowed bool
}

func (x ConditionalExpression) Position() lexer.Position {
	return x.Pos
}

func (x ConditionalExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return Shape{}, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// typecheck pattern
	patternOutputShape, restType, matcher, err := x.Pattern.Typecheck(inputShape)
	if err != nil {
		return Shape{}, nil, err
	}
	// typecheck guard
	var guardOutputShape Shape
	var guardAction states.Action
	if x.Guard == nil {
		guardOutputShape = patternOutputShape
		guardAction = states.SimpleAction(values.BoolValue(true))
	} else {
		guardOutputShape, guardAction, err = x.Guard.Typecheck(patternOutputShape, nil)
		if err != nil {
			return Shape{}, nil, err
		}
		if !(types.BoolType{}).Subsumes(guardOutputShape.Type) {
			return Shape{}, nil, errors.E(
				errors.Code(errors.ConditionMustBeBool),
				errors.Pos(x.Guard.Position()),
				errors.WantType(types.BoolType{}),
				errors.GotType(guardOutputShape.Type),
			)
		}
	}
	// build consequent input shape
	consequentInputShape := Shape{
		Type:  patternOutputShape.Type,
		Stack: guardOutputShape.Stack,
	}
	// typecheck consequent
	consequentOutputShape, consequentAction, err := x.Consequent.Typecheck(consequentInputShape, nil)
	if err != nil {
		return Shape{}, nil, err
	}
	// update input shape
	if x.Guard != nil {
		restType = inputShape.Type
	}
	inputShape = Shape{
		Type:  restType,
		Stack: inputShape.Stack,
	}
	// initialize output type
	outputType := consequentOutputShape.Type
	// typecheck elis patterns, guards, consequents
	elisMatchers := make([]Matcher, len(x.AlternativePatterns))
	elisGuardActions := make([]states.Action, len(x.AlternativeGuards))
	elisConsequentActions := make([]states.Action, len(x.AlternativeConsequents))
	for i := range x.AlternativePatterns {
		// reachability check
		if (types.VoidType{}).Subsumes(inputShape.Type) {
			return Shape{}, nil, errors.E(
				errors.Code(errors.UnreachableElisClause),
				errors.Pos(x.Pattern.Position()),
			)
		}
		// typecheck pattern
		patternOutputShape, restType, elisMatchers[i], err = x.AlternativePatterns[i].Typecheck(inputShape)
		if err != nil {
			return Shape{}, nil, err
		}
		// typecheck guard
		var guardOutputShape Shape
		if x.AlternativeGuards[i] == nil {
			guardOutputShape = patternOutputShape
			elisGuardActions[i] = states.SimpleAction(values.BoolValue(true))
		} else {
			guardOutputShape, elisGuardActions[i], err = x.AlternativeGuards[i].Typecheck(patternOutputShape, nil)
			if err != nil {
				return Shape{}, nil, err
			}
			if !(types.BoolType{}).Subsumes(guardOutputShape.Type) {
				return Shape{}, nil, errors.E(
					errors.Code(errors.ConditionMustBeBool),
					errors.Pos(x.AlternativeGuards[i].Position()),
					errors.WantType(types.BoolType{}),
					errors.GotType(guardOutputShape.Type),
				)
			}
		}
		// build consequent input shape
		consequentInputShape := Shape{
			Type:  patternOutputShape.Type,
			Stack: guardOutputShape.Stack,
		}
		// typecheck consequent
		consequentOutputShape, consequentAction, err := x.AlternativeConsequents[i].Typecheck(consequentInputShape, nil)
		if err != nil {
			return Shape{}, nil, err
		}
		elisConsequentActions[i] = consequentAction
		// update input shape
		if x.AlternativeGuards[i] != nil {
			restType = inputShape.Type
		}
		inputShape = Shape{
			Type:  restType,
			Stack: inputShape.Stack,
		}
		// update output type
		outputType = types.Union(outputType, consequentOutputShape.Type)
	}
	// typecheck alternative
	var alternativeAction states.Action
	if x.Alternative == nil {
		// exhaustivity check
		if !(types.VoidType{}).Subsumes(inputShape.Type) {
			return Shape{}, nil, errors.E(
				errors.Code(errors.NonExhaustiveMatch),
				errors.Pos(x.Pos),
				errors.WantType(types.VoidType{}),
				errors.GotType(inputShape.Type),
			)
		}
	} else {
		// reachability check
		if !x.UnreachableAlternativeAllowed && (types.VoidType{}).Subsumes(inputShape.Type) {
			return Shape{}, nil, errors.E(
				errors.Code(errors.UnreachableElseClause),
				errors.Pos(x.Alternative.Position()),
			)
		}
		// alternative
		var alternativeOutputShape Shape
		alternativeOutputShape, alternativeAction, err = x.Alternative.Typecheck(inputShape, nil)
		if err != nil {
			return Shape{}, nil, err
		}
		// update output type
		outputType = types.Union(outputType, alternativeOutputShape.Type)
	}
	// make action
	action := func(inputState states.State, args []states.Action) states.Thunk {
		matcherVarStack, ok, err := matcher(inputState)
		if err != nil {
			return states.Thunk{State: states.State{}, Drop: false, Err: err}
		}
		if ok {
			guardInputState := states.State{
				Value:     inputState.Value,
				Stack:     matcherVarStack,
				TypeStack: inputState.TypeStack,
			}
			guardState, _, err := guardAction(guardInputState, nil).Eval()
			if err != nil {
				return states.Thunk{State: states.State{}, Drop: false, Err: err}
			}
			boolGuardValue := guardState.Value.(values.BoolValue)
			if boolGuardValue {
				consequentInputState := states.State{
					Value:     inputState.Value,
					Stack:     guardState.Stack,
					TypeStack: inputState.TypeStack,
				}
				return replaceRejectError(consequentAction(consequentInputState, nil), x.Pos)
			}
		}
		for i := range elisMatchers {
			matcherVarStack, ok, err := elisMatchers[i](inputState)
			if err != nil {
				return states.Thunk{State: states.State{}, Drop: false, Err: err}
			}
			if ok {
				guardInputState := states.State{
					Value:     inputState.Value,
					Stack:     matcherVarStack,
					TypeStack: inputState.TypeStack,
				}
				guardState, _, err := elisGuardActions[i](guardInputState, nil).Eval()
				if err != nil {
					return states.Thunk{State: states.State{}, Drop: false, Err: err}
				}
				boolGuardValue := guardState.Value.(values.BoolValue)
				if boolGuardValue {
					consequentInputState := states.State{
						Value:     inputState.Value,
						Stack:     guardState.Stack,
						TypeStack: inputState.TypeStack,
					}
					return replaceRejectError(elisConsequentActions[i](consequentInputState, nil), x.Pos)
				}
			}
		}
		return replaceRejectError(alternativeAction(inputState, nil), x.Pos)
	}
	// return
	outputShape := Shape{
		outputType,
		inputShape.Stack,
	}
	return outputShape, action, nil
}

// replaceRejectError replaces a RejectError with an explainable error with
// location information about the conditional, and the value that was not
// handled.
func replaceRejectError(thunk states.Thunk, pos lexer.Position) states.Thunk {
	if thunk.Func == nil {
		if thunk.Err == nil {
			return thunk
		}
		if rejectError, ok := thunk.Err.(RejectError); ok {
			thunk.Err = errors.E(
				errors.Pos(pos),
				errors.Code(errors.UnexpectedValue),
				errors.GotValue(rejectError.Value),
			)
		}
		return thunk
	}
	return states.Thunk{
		Func: func() states.Thunk {
			return replaceRejectError(thunk.Func(), pos)
		},
	}
}

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
	Name            *string
}

func (p ArrPattern) Position() lexer.Position {
	return p.Pos
}

// spreadInputType spreads the input type for an array pattern over its
// elements and rest.
func spreadInputType(inputType types.Type, elementTypes []types.Type) (restType types.Type, ok bool) {
	switch t := inputType.(type) {
	case *types.NearrType:
		if len(elementTypes) == 0 {
			return t, true
		}
		elementTypes[0] = t.HeadType
		return spreadInputType(t.TailType, elementTypes[1:])
	case *types.ArrType:
		// Optional: fail if the pattern wants to match more elements
		// then the value can contain, as per its type. For now, it is
		// is commented out and will instead lead to an error message
		// about a surplus element not having type Void. That's a bit
		// opaque but has the advantage of indicating the place where
		// the array pattern is too long.
		//if (types.VoidType{}).Subsumes(t.ElType) && len(elementTypes) > 0 {
		//	return nil, false
		//}
		for i := range elementTypes {
			elementTypes[i] = t.ElType
		}
		return t, true
	case types.UnionType:
		for i := range elementTypes {
			elementTypes[i] = types.VoidType{}
		}
		var restType types.Type = types.VoidType{}
		anyOk := false
		for _, disjunct := range t {
			disjunctElementTypes := make([]types.Type, len(elementTypes))
			disjunctRestType, ok := spreadInputType(disjunct, disjunctElementTypes)
			if ok {
				for i := range elementTypes {
					elementTypes[i] = types.Union(elementTypes[i], disjunctElementTypes[i])
				}
				restType = types.Union(restType, disjunctRestType)
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
		return Shape{}, nil, nil, errors.E(
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
	pType := types.NewNearrType(elementTypes, restShape.Type)
	// partition the input type and check for impossible match
	intersection, complement := inputShape.Type.Partition(pType)
	if (types.VoidType{}).Subsumes(intersection) {
		return Shape{}, nil, nil, errors.E(
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
	if p.Name != nil {
		outputShape.Stack = outputShape.Stack.Push(
			VariableFuncer(p, *p.Name, outputShape.Type),
		)
	}
	// build matcher
	matcher := func(inputState states.State) (*states.VariableStack, bool, error) {
		varStack := inputState.Stack
		switch v := inputState.Value.(type) {
		case *values.ArrValue:
			for _, elMatcher := range elementMatchers {
				err := v.Eval()
				if err != nil {
					return nil, false, err
				}
				if v.Head == nil {
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
				v = v.Tail
			}
			varStack, ok, err = restMatcher(states.State{
				Value:     v,
				Stack:     varStack,
				TypeStack: inputState.TypeStack,
			})
			if !ok {
				return nil, false, err
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
	Name           *string
}

func (p ObjPattern) Position() lexer.Position {
	return p.Pos
}

func (p ObjPattern) Typecheck(inputShape Shape) (Shape, types.Type, Matcher, error) {
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
	pType := types.NewObjType(propTypeMap)
	// partition the input type and check for impossible match
	intersection, complement := inputShape.Type.Partition(pType)
	if (types.VoidType{}).Subsumes(intersection) {
		return Shape{}, nil, nil, errors.E(
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
	if p.Name != nil {
		outputShape.Stack = &FuncerStack{
			Head: VariableFuncer(p, *p.Name, outputShape.Type),
			Tail: outputShape.Stack,
		}
	}
	// build matcher
	matcher := func(inputState states.State) (*states.VariableStack, bool, error) {
		varStack := inputState.Stack
		switch v := inputState.Value.(type) {
		case values.ObjValue:
			for prop, valMatcher := range propMatcherMap {
				value, ok := v[prop]
				if !ok {
					return nil, false, nil
				}
				var err error
				varStack, ok, err = valMatcher(states.State{
					Value:     value,
					Stack:     varStack,
					TypeStack: inputState.TypeStack,
				})
				if !ok {
					return nil, false, err
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
	if (types.VoidType{}).Subsumes(intersection) {
		return Shape{}, nil, nil, errors.E(
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
		outputShape.Stack = &FuncerStack{
			Head: VariableFuncer(p, *p.Name, outputShape.Type),
			Tail: outputShape.Stack,
		}
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
