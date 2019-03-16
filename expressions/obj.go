package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/functions"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type ObjExpression struct {
	Pos        lexer.Position
	PropValMap map[string]Expression
}

func (x ObjExpression) Typecheck(inputShape functions.Shape, params []*functions.Parameter) (functions.Shape, functions.Action, error) {
	if len(params) > 0 {
		return zeroShape, nil, errors.E(
			errors.Kind(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	keyTypeMap := make(map[string]types.Type)
	keyActionMap := make(map[string]functions.Action)
	for key, valExpression := range x.PropValMap {
		keyOutputShape, keyAction, err := valExpression.Typecheck(inputShape, nil)
		if err != nil {
			return zeroShape, nil, err
		}
		keyTypeMap[key] = keyOutputShape.Type
		keyActionMap[key] = keyAction
	}
	outputShape := functions.Shape{
		Type:        types.ObjType(keyTypeMap),
		FuncerStack: inputShape.FuncerStack,
	}
	action := func(inputState functions.State, args []functions.Action) functions.State {
		propValMap := make(map[string]values.Value)
		for key, valAction := range keyActionMap {
			valValue := valAction(inputState, nil).Value
			propValMap[key] = valValue
		}
		return functions.State{
			Value: values.ObjValue(propValMap),
			Stack: inputState.Stack,
		}
	}
	return outputShape, action, nil
}
