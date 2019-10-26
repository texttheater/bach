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
	action := func(inputState states.State, args []states.Action) states.State {
		matcherVarStack, ok := matcher(inputState)
		if ok {
			guardInputState := states.State{
				Value:     inputState.Value,
				Stack:     matcherVarStack,
				TypeStack: inputState.TypeStack,
			}
			guardState := guardAction(guardInputState, nil)
			boolGuardValue := guardState.Value.(values.BoolValue)
			if boolGuardValue {
				consequentInputState := states.State{
					Value:     inputState.Value,
					Stack:     guardState.Stack,
					TypeStack: inputState.TypeStack,
				}
				consequentOutputState := consequentAction(consequentInputState, nil)
				return states.State{
					Error:     wrap(consequentOutputState.Error, x.Pos),
					Drop:      consequentOutputState.Drop,
					Value:     consequentOutputState.Value,
					Stack:     inputState.Stack,
					TypeStack: inputState.TypeStack,
				}
			}
		}
		for i := range elisMatchers {
			matcherVarStack, ok := elisMatchers[i](inputState)
			if ok {
				guardInputState := states.State{
					Value:     inputState.Value,
					Stack:     matcherVarStack,
					TypeStack: inputState.TypeStack,
				}
				guardState := elisGuardActions[i](guardInputState, nil)
				boolGuardValue := guardState.Value.(values.BoolValue)
				if boolGuardValue {
					consequentInputState := states.State{
						Value:     inputState.Value,
						Stack:     guardState.Stack,
						TypeStack: inputState.TypeStack,
					}
					consequentOutputState := elisConsequentActions[i](consequentInputState, nil)
					return states.State{
						Error:     wrap(consequentOutputState.Error, x.Pos),
						Drop:      consequentOutputState.Drop,
						Value:     consequentOutputState.Value,
						Stack:     inputState.Stack,
						TypeStack: inputState.TypeStack,
					}
				}
			}
		}
		alternativeOutputState := alternativeAction(inputState, nil)
		return states.State{
			Error:     wrap(alternativeOutputState.Error, x.Pos),
			Drop:      alternativeOutputState.Drop,
			Value:     alternativeOutputState.Value,
			Stack:     inputState.Stack,
			TypeStack: inputState.TypeStack,
		}
	}
	// return
	outputShape := Shape{
		outputType,
		inputShape.Stack,
	}
	return outputShape, action, nil
}

// wrap replaces a RejectError with an explainable error with location
// information about the conditional, and the value that was not handled.
func wrap(e error, pos lexer.Position) error {
	rejectError, ok := e.(RejectError)
	if !ok {
		return e
	}
	return errors.E(
		errors.Pos(pos),
		errors.Code(errors.UnexpectedValue),
		errors.GotValue(rejectError.Value),
	)
}

// pattern/matcher kinda analogous to expression/action

type Pattern interface {
	Position() lexer.Position
	Typecheck(inputShape Shape) (outputShape Shape, restType types.Type, matcher Matcher, err error)
}

type Matcher func(states.State) (*states.VariableStack, bool)

type ArrPattern struct {
	Pos             lexer.Position
	ElementPatterns []Pattern
	Name            *string
}

func (p ArrPattern) Position() lexer.Position {
	return p.Pos
}

func (p ArrPattern) Typecheck(inputShape Shape) (Shape, types.Type, Matcher, error) {
	// compute element input types
	elementInputTypes := make([]types.Type, len(p.ElementPatterns))
	switch t := inputShape.Type.(type) {
	case types.TupType:
		if len(t) != len(elementInputTypes) {
			return Shape{}, types.VoidType{}, nil, nil
		}
		for i, elType := range t {
			elementInputTypes[i] = elType
		}
	case *types.ArrType:
		for i := range elementInputTypes {
			elementInputTypes[i] = t.ElType
		}
	case types.UnionType:
		for i := range elementInputTypes {
			elementInputTypes[i] = types.VoidType{}
		}
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
				return Shape{}, types.VoidType{}, nil, nil
			}
		}
	default:
		return Shape{}, types.VoidType{}, nil, nil
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
	// determine the type of values this pattern will match
	pType := types.TupType(elementTypes)
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
					Value:     v[i],
					Stack:     varStack,
					TypeStack: inputState.TypeStack,
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
					Value:     value,
					Stack:     varStack,
					TypeStack: inputState.TypeStack,
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
	matcher := func(inputState states.State) (*states.VariableStack, bool) {
		// TODO For efficiency, we should check inhabitation of a more
		// general type than p.Type if that is equivalent.
		if inputState.Value.Inhabits(p.Type, inputState.TypeStack) {
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
