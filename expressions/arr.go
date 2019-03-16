package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type ArrExpression struct {
	Pos      lexer.Position
	Elements []Expression
}

func (x ArrExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, states.Action, error) {
	if len(params) > 0 {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	var elementType types.Type = types.VoidType
	elementActions := make([]states.Action, len(x.Elements))
	for i, elExpression := range x.Elements {
		elOutputShape, elAction, err := elExpression.Typecheck(inputShape, nil)
		if err != nil {
			return zeroShape, nil, err
		}
		elementType = types.Union(elementType, elOutputShape.Type)
		elementActions[i] = elAction
	}
	outputShape := functions.Shape{
		Type:        types.ArrType(elementType),
		FuncerStack: inputShape.FuncerStack,
	}
	action := func(inputState states.State, args []states.Action) states.State {
		elementValues := make([]values.Value, len(elementActions))
		for i, elAction := range elementActions {
			elValue := elAction(inputState, nil).Value
			elementValues[i] = elValue
		}
		return states.State{
			Value: values.ArrValue(elementValues),
			Stack: inputState.Stack,
		}
	}
	return outputShape, action, nil
}
