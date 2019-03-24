package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type ObjExpression struct {
	Pos        lexer.Position
	PropValMap map[string]Expression
}

func (x ObjExpression) Typecheck(inputShape shapes.Shape, params []*shapes.Parameter) (shapes.Shape, states.Action, error) {
	if len(params) > 0 {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	keyTypeMap := make(map[string]types.Type)
	keyActionMap := make(map[string]states.Action)
	for key, valExpression := range x.PropValMap {
		keyOutputShape, keyAction, err := valExpression.Typecheck(inputShape, nil)
		if err != nil {
			return zeroShape, nil, err
		}
		keyTypeMap[key] = keyOutputShape.Type
		keyActionMap[key] = keyAction
	}
	outputShape := shapes.Shape{
		Type:        types.ObjType(keyTypeMap),
		FuncerStack: inputShape.FuncerStack,
	}
	action := func(inputState states.State, args []states.Action) states.State {
		propValMap := make(map[string]values.Value)
		for key, valAction := range keyActionMap {
			valValue := valAction(inputState, nil).Value
			propValMap[key] = valValue
		}
		return states.State{
			Value: values.ObjValue(propValMap),
			Stack: inputState.Stack,
		}
	}
	return outputShape, action, nil
}
