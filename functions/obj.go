package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type ObjExpression struct {
	Pos        lexer.Position
	PropValMap map[string]Expression
}

func (x ObjExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	if len(params) > 0 {
		return Shape{}, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	keyTypeMap := make(map[string]types.Type)
	keyActionMap := make(map[string]states.Action)
	for key, valExpression := range x.PropValMap {
		keyOutputShape, keyAction, err := valExpression.Typecheck(inputShape, nil)
		if err != nil {
			return Shape{}, nil, err
		}
		keyTypeMap[key] = keyOutputShape.Type
		keyActionMap[key] = keyAction
	}
	outputShape := Shape{
		Type:  types.NewObjType(keyTypeMap),
		Stack: inputShape.Stack,
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