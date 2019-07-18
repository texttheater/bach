package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/patterns"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type MatchExpression struct {
	Pos             lexer.Position
	Pattern         patterns.Pattern
	Guard           Expression
	Consequent      Expression
	ElisPatterns    []patterns.Pattern
	ElisGuards      []Expression
	ElisConsequents []Expression
	Alternative     Expression
}

func (x MatchExpression) Typecheck(inputShape shapes.Shape, params []*shapes.Parameter) (shapes.Shape, states.Action, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return shapes.Shape{}, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// typecheck pattern
	patternOutputShape, restType, matcher, err := x.Pattern.Typecheck(inputShape)
	if err != nil {
		return shapes.Shape{}, nil, err
	}
	// typecheck guard
	var guardOutputShape shapes.Shape
	var guardAction states.Action
	if x.Guard == nil {
		guardOutputShape = patternOutputShape
		guardAction = states.SimpleAction(values.BoolValue(true))
	} else {
		guardOutputShape, guardAction, err = x.Guard.Typecheck(patternOutputShape, nil)
		if err != nil {
			return shapes.Shape{}, nil, err
		}
		if !(types.BoolType{}).Subsumes(guardOutputShape.Type) {
			return shapes.Shape{}, nil, errors.E(
				errors.Code(errors.ConditionMustBeBool),
				errors.Pos(x.Pos),
				errors.WantType(types.BoolType{}),
				errors.GotType(guardOutputShape.Type),
			)
		}
	}
	// build consequent input shape
	consequentInputShape := shapes.Shape{
		Type:  patternOutputShape.Type,
		Stack: guardOutputShape.Stack,
	}
	// typecheck consequent
	consequentOutputShape, consequentAction, err := x.Consequent.Typecheck(consequentInputShape, nil)
	if err != nil {
		return shapes.Shape{}, nil, err
	}
	// update input shape
	if x.Guard != nil {
		restType = inputShape.Type
	}
	inputShape = shapes.Shape{
		Type:  restType,
		Stack: inputShape.Stack,
	}
	// initialize output type
	outputType := consequentOutputShape.Type
	// typecheck elis patterns, guards, consequents
	elisMatchers := make([]patterns.Matcher, len(x.ElisPatterns))
	elisGuardActions := make([]states.Action, len(x.ElisGuards))
	elisConsequentActions := make([]states.Action, len(x.ElisConsequents))
	for i := range x.ElisPatterns {
		// reachability check
		if (types.VoidType{}).Subsumes(inputShape.Type) {
			return shapes.Shape{}, nil, errors.E(
				errors.Code(errors.UnreachableElisClause),
				errors.Pos(x.Pos),
			)
		}
		// typecheck pattern
		patternOutputShape, restType, elisMatchers[i], err = x.ElisPatterns[i].Typecheck(inputShape)
		if err != nil {
			return shapes.Shape{}, nil, err
		}
		// typecheck guard
		var guardOutputShape shapes.Shape
		if x.ElisGuards[i] == nil {
			guardOutputShape = patternOutputShape
			elisGuardActions[i] = states.SimpleAction(values.BoolValue(true))
		} else {
			guardOutputShape, elisGuardActions[i], err = x.ElisGuards[i].Typecheck(patternOutputShape, nil)
			if err != nil {
				return shapes.Shape{}, nil, err
			}
			if !(types.BoolType{}).Subsumes(guardOutputShape.Type) {
				return shapes.Shape{}, nil, errors.E(
					errors.Code(errors.ConditionMustBeBool),
					errors.Pos(x.Pos),
					errors.WantType(types.BoolType{}),
					errors.GotType(guardOutputShape.Type),
				)
			}
		}
		// build consequent input shape
		consequentInputShape := shapes.Shape{
			Type:  patternOutputShape.Type,
			Stack: guardOutputShape.Stack,
		}
		// typecheck consequent
		consequentOutputShape, consequentAction, err := x.ElisConsequents[i].Typecheck(consequentInputShape, nil)
		if err != nil {
			return shapes.Shape{}, nil, err
		}
		elisConsequentActions[i] = consequentAction
		// update input shape
		if x.ElisGuards[i] != nil {
			restType = inputShape.Type
		}
		inputShape = shapes.Shape{
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
			return shapes.Shape{}, nil, errors.E(
				errors.Code(errors.NonExhaustiveMatch),
				errors.Pos(x.Pos),
				errors.WantType(types.VoidType{}),
				errors.GotType(inputShape.Type),
			)
		}
	} else {
		// reachability check
		if (types.VoidType{}).Subsumes(inputShape.Type) {
			return shapes.Shape{}, nil, errors.E(
				errors.Code(errors.UnreachableElseClause),
				errors.Pos(x.Pos),
			)
		}
		// alternative
		var alternativeOutputShape shapes.Shape
		alternativeOutputShape, alternativeAction, err = x.Alternative.Typecheck(inputShape, nil)
		if err != nil {
			return shapes.Shape{}, nil, err
		}
		outputType = types.Union(outputType, alternativeOutputShape.Type)
	}
	// make action
	action := func(inputState states.State, args []states.Action) states.State {
		matcherVarStack, ok := matcher(inputState)
		if ok {
			guardInputState := states.State{
				Value: inputState.Value,
				Stack: matcherVarStack,
			}
			guardState := guardAction(guardInputState, nil)
			boolGuardValue := guardState.Value.(values.BoolValue)
			if boolGuardValue {
				consequentInputState := states.State{
					Value: inputState.Value,
					Stack: guardState.Stack,
				}
				consequentOutputState := consequentAction(consequentInputState, nil)
				return states.State{
					Value: consequentOutputState.Value,
					Stack: inputState.Stack,
				}
			}
		}
		for i := range elisMatchers {
			matcherVarStack, ok := elisMatchers[i](inputState)
			if ok {
				guardInputState := states.State{
					Value: inputState.Value,
					Stack: matcherVarStack,
				}
				guardState := elisGuardActions[i](guardInputState, nil)
				boolGuardValue := guardState.Value.(values.BoolValue)
				if boolGuardValue {
					consequentInputState := states.State{
						Value: inputState.Value,
						Stack: guardState.Stack,
					}
					consequentOutputState := elisConsequentActions[i](consequentInputState, nil)
					return states.State{
						Value: consequentOutputState.Value,
						Stack: inputState.Stack,
					}
				}
			}
		}
		alternativeOutputState := alternativeAction(inputState, nil)
		return states.State{
			alternativeOutputState.Value,
			inputState.Stack,
		}
	}
	// return
	outputShape := shapes.Shape{
		outputType,
		inputShape.Stack,
	}
	return outputShape, action, nil
}
