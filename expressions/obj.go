package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/parameters"
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

func (x ObjExpression) Typecheck(inputShape Shape, params []*parameters.Parameter) (Shape, states.Action, *states.IDStack, error) {
	if len(params) > 0 {
		return Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	propTypeMap := make(map[string]types.Type)
	propActionMap := make(map[string]states.Action)
	var ids *states.IDStack
	for prop, valExpression := range x.PropValMap {
		propOutputShape, propAction, propIDs, err := valExpression.Typecheck(inputShape, nil)
		if err != nil {
			return Shape{}, nil, nil, err
		}
		propTypeMap[prop] = propOutputShape.Type
		propActionMap[prop] = propAction
		ids = ids.AddAll(propIDs)
	}
	outputShape := Shape{
		Type: types.Obj{
			Props: propTypeMap,
			Rest:  types.Void{},
		},
		Stack: inputShape.Stack,
	}
	action := func(inputState states.State, args []states.Action) *states.Thunk {
		propThunkMap := make(map[string]*states.Thunk)
		for prop, valAction := range propActionMap {
			propThunkMap[prop] = valAction(inputState, nil)
		}
		return states.ThunkFromState(states.State{
			Value:     states.ObjValue(propThunkMap),
			Stack:     inputState.Stack,
			TypeStack: inputState.TypeStack,
		})
	}
	return outputShape, action, ids, nil
}
