package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type AssignmentExpression struct {
	Pos     lexer.Position
	Pattern Pattern
}

func (x AssignmentExpression) Position() lexer.Position {
	return x.Pos
}

func (x AssignmentExpression) Typecheck(inputShape Shape, params []*params.Param) (Shape, states.Action, *states.IDStack, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos))

	}
	// typecheck pattern
	patternOutputShape, restType, matcher, err := x.Pattern.Typecheck(inputShape)
	if err != nil {
		return Shape{}, nil, nil, err
	}
	// make sure pattern matches input Type
	if !(types.Void{}).Subsumes(restType) {
		return Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.NonExhaustiveMatch),
			errors.Pos(x.Pos),
			errors.GotType(types.Void{}),
			errors.GotType(restType),
		)
	}
	// make action
	action := func(inputState states.State, args []states.Action) *states.Thunk {
		matcherVarStack, _, err := matcher(inputState)
		if err != nil {
			return states.ThunkFromError(err)
		}
		return states.ThunkFromState(states.State{
			Value:     inputState.Value,
			Stack:     matcherVarStack,
			TypeStack: inputState.TypeStack,
		})
	}
	return patternOutputShape, action, nil, nil
}
