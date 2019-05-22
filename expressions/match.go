package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/patterns"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type MatchExpression struct {
	Pos             lexer.Position
	Pattern         patterns.Pattern
	Consequent      Expression
	ElisPatterns    []patterns.Pattern
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
	consequentInputShape, restType, matcher, err := x.Pattern.Typecheck(inputShape)
	if err != nil {
		return shapes.Shape{}, nil, err
	}
	// typecheck consequent
	consequentOutputShape, consequentAction, err := x.Consequent.Typecheck(consequentInputShape, nil)
	if err != nil {
		return shapes.Shape{}, nil, err
	}
	// update input shape
	inputShape = shapes.Shape{
		Type:  restType,
		Stack: inputShape.Stack,
	}
	// initialize output type
	outputType := consequentOutputShape.Type
	// typecheck elis patterns, consequents
	elisMatchers := make([]patterns.Matcher, len(x.ElisPatterns))
	elisConsequentActions := make([]states.Action, len(x.ElisConsequents))
	for i := range x.ElisPatterns {
		consequentInputShape, restType, elisMatchers[i], err = x.ElisPatterns[i].Typecheck(inputShape)
		if err != nil {
			return shapes.Shape{}, nil, err
		}
		consequentOutputShape, consequentAction, err := x.ElisConsequents[i].Typecheck(consequentInputShape, nil)
		if err != nil {
			return shapes.Shape{}, nil, err
		}
		elisConsequentActions[i] = consequentAction
		// update input shape
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
		if matcherVarStack, ok := matcher(inputState); ok {
			consequentInputState := states.State{
				Value: inputState.Value,
				Stack: matcherVarStack,
			}
			consequentOutputState := consequentAction(consequentInputState, nil)
			return states.State{
				consequentOutputState.Value,
				inputState.Stack,
			}
		}
		for i := range elisMatchers {
			if matcherVarStack, ok := elisMatchers[i](inputState); ok {
				consequentInputState := states.State{
					Value: inputState.Value,
					Stack: matcherVarStack,
				}
				consequentOutputState := elisConsequentActions[i](consequentInputState, nil)
				return states.State{
					consequentOutputState.Value,
					inputState.Stack,
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
