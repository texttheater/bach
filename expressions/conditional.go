package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
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

func (x ConditionalExpression) Typecheck(inputShape shapes.Shape, params []*params.Param) (shapes.Shape, states.Action, *states.IDStack, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return shapes.Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// typecheck pattern
	patternOutputShape, restType, matcher, err := x.Pattern.Typecheck(inputShape)
	if err != nil {
		return shapes.Shape{}, nil, nil, err
	}
	// typecheck guard
	var guardOutputShape shapes.Shape
	var guardAction states.Action
	var guardIDs *states.IDStack
	if x.Guard == nil {
		guardOutputShape = patternOutputShape
		guardAction = states.SimpleAction(states.BoolValue(true))
	} else {
		guardOutputShape, guardAction, guardIDs, err = x.Guard.Typecheck(patternOutputShape, nil)
		if err != nil {
			return shapes.Shape{}, nil, nil, err
		}
		if !(types.BoolType{}).Subsumes(guardOutputShape.Type) {
			return shapes.Shape{}, nil, nil, errors.TypeError(
				errors.Code(errors.ConditionMustBeBool),
				errors.Pos(x.Guard.Position()),
				errors.WantType(types.BoolType{}),
				errors.GotType(guardOutputShape.Type),
			)
		}
	}
	ids := guardIDs
	// build consequent input shape
	consequentInputShape := shapes.Shape{
		Type:  patternOutputShape.Type,
		Stack: guardOutputShape.Stack,
	}
	// typecheck consequent
	consequentOutputShape, consequentAction, consequentIDs, err := x.Consequent.Typecheck(consequentInputShape, nil)
	if err != nil {
		return shapes.Shape{}, nil, nil, err
	}
	ids = ids.AddAll(consequentIDs)
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
	elisMatchers := make([]Matcher, len(x.AlternativePatterns))
	elisGuardActions := make([]states.Action, len(x.AlternativeGuards))
	elisConsequentActions := make([]states.Action, len(x.AlternativeConsequents))
	for i := range x.AlternativePatterns {
		// reachability check
		if (types.VoidType{}).Subsumes(inputShape.Type) {
			return shapes.Shape{}, nil, nil, errors.TypeError(
				errors.Code(errors.UnreachableElisClause),
				errors.Pos(x.Pattern.Position()),
			)
		}
		// typecheck pattern
		patternOutputShape, restType, elisMatchers[i], err = x.AlternativePatterns[i].Typecheck(inputShape)
		if err != nil {
			return shapes.Shape{}, nil, nil, err
		}
		// typecheck guard
		var guardOutputShape shapes.Shape
		if x.AlternativeGuards[i] == nil {
			guardOutputShape = patternOutputShape
			elisGuardActions[i] = states.SimpleAction(states.BoolValue(true))
		} else {
			guardOutputShape, elisGuardActions[i], guardIDs, err = x.AlternativeGuards[i].Typecheck(patternOutputShape, nil)
			if err != nil {
				return shapes.Shape{}, nil, nil, err
			}
			if !(types.BoolType{}).Subsumes(guardOutputShape.Type) {
				return shapes.Shape{}, nil, nil, errors.TypeError(
					errors.Code(errors.ConditionMustBeBool),
					errors.Pos(x.AlternativeGuards[i].Position()),
					errors.WantType(types.BoolType{}),
					errors.GotType(guardOutputShape.Type),
				)
			}
			ids = ids.AddAll(guardIDs)
		}
		// build consequent input shape
		consequentInputShape := shapes.Shape{
			Type:  patternOutputShape.Type,
			Stack: guardOutputShape.Stack,
		}
		// typecheck consequent
		consequentOutputShape, consequentAction, consequentIDs, err := x.AlternativeConsequents[i].Typecheck(consequentInputShape, nil)
		if err != nil {
			return shapes.Shape{}, nil, nil, err
		}
		elisConsequentActions[i] = consequentAction
		ids = ids.AddAll(consequentIDs)
		// update input shape
		if x.AlternativeGuards[i] != nil {
			restType = inputShape.Type
		}
		inputShape = shapes.Shape{
			Type:  restType,
			Stack: inputShape.Stack,
		}
		// update output type
		outputType = types.NewUnionType(outputType, consequentOutputShape.Type)
	}
	// typecheck alternative
	var alternativeAction states.Action
	var alternativeIDs *states.IDStack
	if x.Alternative == nil {
		// exhaustivity check
		//if !(types.VoidTypeType{}).Subsumes(inputShape.Type) {
		//	return shapes.Shape{}, nil, nil, errors.TypeError(
		//		errors.Code(errors.NonExhaustiveMatch),
		//		errors.Pos(x.Pos),
		//		errors.WantType(types.VoidTypeType{}),
		//		errors.GotType(inputShape.Type),
		//	)
		//}
	} else {
		// reachability check
		if !x.UnreachableAlternativeAllowed && (types.VoidType{}).Subsumes(inputShape.Type) {
			return shapes.Shape{}, nil, nil, errors.TypeError(
				errors.Code(errors.UnreachableElseClause),
				errors.Pos(x.Alternative.Position()),
			)
		}
		// alternative
		var alternativeOutputShape shapes.Shape
		alternativeOutputShape, alternativeAction, alternativeIDs, err = x.Alternative.Typecheck(inputShape, nil)
		if err != nil {
			return shapes.Shape{}, nil, nil, err
		}
		ids = ids.AddAll(alternativeIDs)
		// update output type
		outputType = types.NewUnionType(outputType, alternativeOutputShape.Type)
	}
	// make action
	action := func(inputState states.State, args []states.Action) *states.Thunk {
		matcherVarStack, ok, err := matcher(inputState)
		if err != nil {
			return states.ThunkFromError(err)
		}
		if ok {
			guardInputState := states.State{
				Value:     inputState.Value,
				Stack:     matcherVarStack,
				TypeStack: inputState.TypeStack,
			}
			thunk := guardAction(guardInputState, nil)
			val, err := thunk.Eval()
			if err != nil {
				return states.ThunkFromError(err)
			}
			boolGuardValue := val.(states.BoolValue)
			if boolGuardValue {
				consequentInputState := states.State{
					Value:     inputState.Value,
					Stack:     thunk.Stack,
					TypeStack: inputState.TypeStack,
				}
				return consequentAction(consequentInputState, nil)
			}
		}
		for i := range elisMatchers {
			matcherVarStack, ok, err := elisMatchers[i](inputState)
			if err != nil {
				return states.ThunkFromError(err)
			}
			if ok {
				guardInputState := states.State{
					Value:     inputState.Value,
					Stack:     matcherVarStack,
					TypeStack: inputState.TypeStack,
				}
				thunk := elisGuardActions[i](guardInputState, nil)
				val, err := thunk.Eval()
				if err != nil {
					return states.ThunkFromError(err)
				}
				boolGuardValue := val.(states.BoolValue)
				if boolGuardValue {
					consequentInputState := states.State{
						Value:     inputState.Value,
						Stack:     thunk.Stack,
						TypeStack: inputState.TypeStack,
					}
					return elisConsequentActions[i](consequentInputState, nil)
				}
			}
		}
		if alternativeAction == nil {
			return states.ThunkFromError(errors.TypeError(
				errors.Pos(x.Pos),
				errors.Code(errors.UnexpectedValue),
				errors.GotValue(inputState.Value)))
		}
		return alternativeAction(inputState, nil)
	}
	// return
	outputShape := shapes.Shape{
		outputType,
		inputShape.Stack,
	}
	return outputShape, action, ids, nil
}
