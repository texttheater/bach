package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type CompositionExpression struct {
	Pos   lexer.Position
	Left  Expression
	Right Expression
}

func (x CompositionExpression) Position() lexer.Position {
	return x.Pos
}

func (x CompositionExpression) Typecheck(inputShape Shape, params []*Parameter) (Shape, states.Action, error) {
	if len(params) > 0 {
		return Shape{}, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	middleShape, lAction, err := x.Left.Typecheck(inputShape, nil)
	if err != nil {
		return Shape{}, nil, err
	}
	if (types.VoidType{}).Subsumes(middleShape.Type) {
		return Shape{}, nil, errors.E(
			errors.Code(errors.ComposeWithVoid),
			errors.Pos(x.Right.Position()),
		)
	}
	outputShape, rAction, err := x.Right.Typecheck(middleShape, nil)
	if err != nil {
		return Shape{}, nil, err
	}
	action := func(inputState states.State, args []states.Action) states.State {
		middleState := lAction(inputState, nil)
		if middleState.Drop {
			return middleState
		}
		outputState := rAction(middleState, nil)
		return outputState
	}
	return outputShape, action, nil
}

func Compose(pos lexer.Position, l Expression, r Expression) Expression {
	if l == nil {
		return r
	}
	return &CompositionExpression{pos, l, r}
}
