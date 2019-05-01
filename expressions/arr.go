package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type ArrExpression struct {
	Pos      lexer.Position
	Elements []Expression
}

func (x ArrExpression) Typecheck(inputShape shapes.Shape, params []*shapes.Parameter) (shapes.Shape, states.Action, error) {
	if len(params) > 0 {
		return shapes.Shape{}, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	arrType := types.VoidArrType
	elementActions := make([]states.Action, len(x.Elements))
	for i := len(x.Elements) - 1; i >= 0; i-- {
		elExpression := x.Elements[i]
		elOutputShape, elAction, err := elExpression.Typecheck(inputShape, nil)
		if err != nil {
			return shapes.Shape{}, nil, err
		}
		arrType = &types.NearrType{elOutputShape.Type, arrType}
		elementActions[i] = elAction
	}
	outputShape := shapes.Shape{
		Type:        arrType,
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
