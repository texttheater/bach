package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type ConditionalExpression struct {
	Pos             lexer.Position
	Condition       Expression
	Consequent      Expression
	ElifConditions  []Expression
	ElifConsequents []Expression
	Alternative     Expression
}

func (x ConditionalExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, states.Action, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// typecheck condition
	conditionOutputShape, conditionAction, err := x.Condition.Typecheck(inputShape, nil)
	if err != nil {
		return zeroShape, nil, err
	}
	if !types.BoolType.Subsumes(conditionOutputShape.Type) {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.ConditionMustBeBool),
			errors.Pos(x.Pos),
			errors.WantType(types.BoolType),
			errors.GotType(conditionOutputShape.Type),
		)
	}
	// context is the shared input context for all conditions and consequents.
	// Each condition may add to the FuncerStack. Type always stays the same.
	shape := functions.Shape{
		Type:        inputShape.Type,
		FuncerStack: conditionOutputShape.FuncerStack,
	}
	consequentOutputShape, consequentAction, err := x.Consequent.Typecheck(shape, nil)
	if err != nil {
		return zeroShape, nil, err
	}
	outputType := consequentOutputShape.Type
	elifConditionActions := make([]states.Action, len(x.ElifConditions))
	elifConsequentActions := make([]states.Action, len(x.ElifConsequents))
	for i := range x.ElifConditions {
		conditionOutputShape, elifConditionAction, err := x.ElifConditions[i].Typecheck(shape, nil)
		if err != nil {
			return zeroShape, nil, err
		}
		if !types.BoolType.Subsumes(conditionOutputShape.Type) {
			return zeroShape, nil, errors.E(
				errors.Kind(errors.ConditionMustBeBool),
				errors.Pos(x.Pos),
				errors.WantType(types.BoolType),
				errors.GotType(conditionOutputShape.Type),
			)
		}
		shape.FuncerStack = conditionOutputShape.FuncerStack
		elifConditionActions[i] = elifConditionAction
		consequentOutputShape, elifConsequentAction, err := x.ElifConsequents[i].Typecheck(shape, nil)
		if err != nil {
			return zeroShape, nil, err
		}
		elifConsequentActions[i] = elifConsequentAction
		outputType = types.Union(outputType, consequentOutputShape.Type)
	}
	alternativeOutputShape, alternativeAction, err := x.Alternative.Typecheck(shape, nil)
	if err != nil {
		return zeroShape, nil, err
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
	outputShape := functions.Shape{
		Type:        outputType,
		FuncerStack: inputShape.FuncerStack,
	}
	return outputShape, action, nil
}
