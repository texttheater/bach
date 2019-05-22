package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

// TODO Merge conditional expressions into match expressions?

type ConditionalExpression struct {
	Pos             lexer.Position
	Condition       Expression
	Consequent      Expression
	ElifConditions  []Expression
	ElifConsequents []Expression
	Alternative     Expression
}

func (x ConditionalExpression) Typecheck(inputShape shapes.Shape, params []*shapes.Parameter) (shapes.Shape, states.Action, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return shapes.Shape{}, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// typecheck condition
	conditionOutputShape, conditionAction, err := x.Condition.Typecheck(inputShape, nil)
	if err != nil {
		return shapes.Shape{}, nil, err
	}
	if !(types.BoolType{}).Subsumes(conditionOutputShape.Type) {
		return shapes.Shape{}, nil, errors.E(
			errors.Code(errors.ConditionMustBeBool),
			errors.Pos(x.Pos),
			errors.WantType(types.BoolType{}),
			errors.GotType(conditionOutputShape.Type),
		)
	}
	// shape is the shared input shape for all conditions and consequents.
	// Each condition may add to the Stack. Type always stays the same.
	shape := shapes.Shape{
		Type:  inputShape.Type,
		Stack: conditionOutputShape.Stack,
	}
	consequentOutputShape, consequentAction, err := x.Consequent.Typecheck(shape, nil)
	if err != nil {
		return shapes.Shape{}, nil, err
	}
	outputType := consequentOutputShape.Type
	elifConditionActions := make([]states.Action, len(x.ElifConditions))
	elifConsequentActions := make([]states.Action, len(x.ElifConsequents))
	for i := range x.ElifConditions {
		conditionOutputShape, elifConditionAction, err := x.ElifConditions[i].Typecheck(shape, nil)
		if err != nil {
			return shapes.Shape{}, nil, err
		}
		if !(types.BoolType{}).Subsumes(conditionOutputShape.Type) {
			return shapes.Shape{}, nil, errors.E(
				errors.Code(errors.ConditionMustBeBool),
				errors.Pos(x.Pos),
				errors.WantType(types.BoolType{}),
				errors.GotType(conditionOutputShape.Type),
			)
		}
		shape.Stack = conditionOutputShape.Stack
		elifConditionActions[i] = elifConditionAction
		consequentOutputShape, elifConsequentAction, err := x.ElifConsequents[i].Typecheck(shape, nil)
		if err != nil {
			return shapes.Shape{}, nil, err
		}
		elifConsequentActions[i] = elifConsequentAction
		outputType = types.Union(outputType, consequentOutputShape.Type)
	}
	alternativeOutputShape, alternativeAction, err := x.Alternative.Typecheck(shape, nil)
	if err != nil {
		return shapes.Shape{}, nil, err
	}
	outputType = types.Union(outputType, alternativeOutputShape.Type)
	action := func(inputState states.State, args []states.Action) states.State {
		conditionState := conditionAction(inputState, nil)
		boolConditionValue, _ := conditionState.Value.(values.BoolValue)
		if boolConditionValue {
			consequentInputState := states.State{
				Value: inputState.Value,
				Stack: conditionState.Stack,
			}
			consequentOutputState := consequentAction(consequentInputState, nil)
			return states.State{
				Value: consequentOutputState.Value,
				Stack: inputState.Stack,
			}
		}
		for i := range elifConditionActions {
			conditionState = elifConditionActions[i](inputState, nil)
			boolConditionValue, _ = conditionState.Value.(values.BoolValue)
			if boolConditionValue {
				consequentInputState := states.State{
					Value: inputState.Value,
					Stack: conditionState.Stack,
				}
				consequentOutputState := elifConsequentActions[i](consequentInputState, nil)
				return states.State{
					Value: consequentOutputState.Value,
					Stack: inputState.Stack,
				}
			}
		}
		alternativeOutputState := alternativeAction(inputState, nil)
		return states.State{
			Value: alternativeOutputState.Value,
			Stack: inputState.Stack,
		}
	}
	outputShape := shapes.Shape{
		Type:  outputType,
		Stack: inputShape.Stack,
	}
	return outputShape, action, nil
}
