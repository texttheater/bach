package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type ObjExpression struct {
	Pos        lexer.Position
	PropValMap map[string]Expression
}

func (x ObjExpression) Position() lexer.Position {
	return x.Pos
}

func (x ObjExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, *states.IDStack, error) {
	if len(params) > 0 {
		return Shape{}, nil, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	keyTypeMap := make(map[string]types.Type)
	keyActionMap := make(map[string]states.Action)
	var ids *states.IDStack
	for key, valExpression := range x.PropValMap {
		keyOutputShape, keyAction, keyIDs, err := valExpression.Typecheck(inputShape, nil)
		if err != nil {
			return Shape{}, nil, nil, err
		}
		keyTypeMap[key] = keyOutputShape.Type
		keyActionMap[key] = keyAction
		ids = ids.AddAll(keyIDs)
	}
	outputShape := Shape{
		Type:  types.NewObjType(keyTypeMap),
		Stack: inputShape.Stack,
	}
	action := func(inputState states.State, args []states.Action) *states.Thunk {
		propThunkMap := make(map[string]*states.Thunk)
		for key, valAction := range keyActionMap {
			propThunkMap[key] = valAction(inputState, nil)
		}
		return states.ThunkFromState(states.State{
			Value:     states.ObjValue(propThunkMap),
			Stack:     inputState.Stack,
			TypeStack: inputState.TypeStack,
		})
	}
	return outputShape, action, ids, nil
}
