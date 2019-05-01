package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/shapes"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type CompositionExpression struct {
	Pos   lexer.Position
	Left  Expression
	Right Expression
}

func (x CompositionExpression) Typecheck(inputShape shapes.Shape, params []*shapes.Parameter) (shapes.Shape, states.Action, error) {
	if len(params) > 0 {
		return shapes.Shape{}, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	middleShape, lAction, err := x.Left.Typecheck(inputShape, nil)
	if err != nil {
		return shapes.Shape{}, nil, err
	}
	if (types.VoidType{}).Subsumes(middleShape.Type) {
		return shapes.Shape{}, nil, errors.E(
			errors.Code(errors.ComposeWithVoid),
			errors.Pos(x.Pos),
		)
	}
	outputShape, rAction, err := x.Right.Typecheck(middleShape, nil)
	if err != nil {
		return shapes.Shape{}, nil, err
	}
	action := func(inputState states.State, args []states.Action) states.State {
		middleState := lAction(inputState, nil)
		outputState := rAction(middleState, nil)
		return outputState
	}
	return outputShape, action, nil
}
