package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
	"github.com/texttheater/bach/values"
)

type RejectExpression struct {
	Pos lexer.Position
}

func (x RejectExpression) Position() lexer.Position {
	return x.Pos
}

func (x RejectExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return Shape{}, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// create output shape
	outputShape := Shape{
		Type: types.VoidType{},
	}
	// create action
	action := func(inputState states.State, args []states.Action) states.State {
		return states.State{
			Error: RejectError{
				Value: inputState.Value,
			},
		}
	}
	// return
	return outputShape, action, nil
}

type RejectError struct {
	Value values.Value
}

func (e RejectError) Error() string {
	return e.Value.String() + ": value rejected"
}
