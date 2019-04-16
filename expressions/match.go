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
		return zeroShape, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// typecheck pattern
	shape, matcher, err := x.Pattern.Typecheck(inputShape)
	if err != nil {
		return zeroShape, nil, err
	}
	// typecheck consequent
	consequentOutputShape, consequentAction, err := x.Consequent.Typecheck(shape, nil)
	if err != nil {
		return zeroShape, nil, err
	}
	// typecheck elis patterns, consequents
	outputType := consequentOutputShape.Type
	elisMatchers := make([]patterns.Matcher, len(x.ElisPatterns))
	elisConsequentActions := make([]states.Action, len(x.ElisConsequents))
	for i := range x.ElisPatterns {
		shape, elisMatchers[i], err = x.ElisPatterns[i].Typecheck(shape)
		if err != nil {
			return zeroShape, nil, err
		}
		consequentOutputShape, consequentAction, err := x.ElisConsequents[i].Typecheck(shape, nil)
		if err != nil {
			return zeroShape, nil, err
		}
		elisConsequentActions[i] = consequentAction
		outputType = types.Union(outputType, consequentOutputShape.Type)
	}
	// typecheck alternative
	var alternativeAction states.Action
	if x.Alternative == nil {
		// exhaustivity check
		if !types.VoidType.Subsumes(shape.Type) {
			return zeroShape, nil, errors.E(
				errors.Code(errors.NonExhaustiveMatch),
				errors.Pos(x.Pos),
				errors.WantType(types.VoidType),
				errors.GotType(shape.Type),
			)
		}
	} else {
		var alternativeOutputShape shapes.Shape
		alternativeOutputShape, alternativeAction, err = x.Alternative.Typecheck(shape, nil)
		if err != nil {
			return zeroShape, nil, err
		}
		outputType = types.Union(outputType, alternativeOutputShape.Type)
	}
	// make action
	action := func(inputState states.State, args []states.Action) states.State {
		if matcherState, ok := matcher(inputState); ok {
			consequentOutputState := consequentAction(matcherState, nil)
			return states.State{
				consequentOutputState.Value,
				inputState.Stack,
			}
		}
		for i := range elisMatchers {
			if matcherState, ok := elisMatchers[i](inputState); ok {
				consequentOutputState := elisConsequentActions[i](matcherState, nil)
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
		inputShape.FuncerStack,
	}
	return outputShape, action, nil
}
