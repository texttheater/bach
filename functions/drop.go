package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type DropExpression struct {
	Pos lexer.Position
}

func (x DropExpression) Position() lexer.Position {
	return x.Pos
}

func (x DropExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return Shape{}, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// create output shape
	outputShape := Shape{
		Type:  types.VoidType{},
		Stack: nil,
	}
	// create action
	action := func(inputState states.State, args []states.Action) states.Thunk {
		return states.EagerThunk(states.State{}, true, nil)
	}
	// return
	return outputShape, action, nil
}
