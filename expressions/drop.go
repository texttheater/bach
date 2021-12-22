package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type DropExpression struct {
	Pos lexer.Position
}

func (x DropExpression) Position() lexer.Position {
	return x.Pos
}

func (x DropExpression) Typecheck(inputShape Shape, params []*parameters.Parameter) (Shape, states.Action, *states.IDStack, error) {
	// make sure we got no parameters
	if len(params) > 0 {
		return Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	// create output shape
	outputShape := Shape{
		Type:  types.Void{},
		Stack: nil,
	}
	// create action
	action := func(inputState states.State, args []states.Action) *states.Thunk {
		return &states.Thunk{
			Result: states.Result{
				Drop: true,
			},
		}
	}
	// return
	return outputShape, action, nil, nil
}
