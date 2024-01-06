package expressions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
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

func (x CompositionExpression) Typecheck(inputShape shapes.Shape, params []*params.Param) (shapes.Shape, states.Action, *states.IDStack, error) {
	if len(params) > 0 {
		return shapes.Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos))

	}
	middleShape, lAction, ids, err := x.Left.Typecheck(inputShape, nil)
	if err != nil {
		return shapes.Shape{}, nil, nil, err
	}
	if (types.Void{}).Subsumes(middleShape.Type) {
		return shapes.Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.ComposeWithVoid),
			errors.Pos(x.Right.Position()))

	}
	outputShape, rAction, rIDs, err := x.Right.Typecheck(middleShape, nil)
	if err != nil {
		return shapes.Shape{}, nil, nil, err
	}
	action := func(inputState states.State, args []states.Action) *states.Thunk {
		thunk := lAction(inputState, nil)
		val, err := thunk.Eval()
		if err != nil {
			return states.ThunkFromError(err)
		}
		state := states.State{
			Value:     val,
			Stack:     thunk.Stack,
			TypeStack: thunk.TypeStack,
		}
		return rAction(state, nil)
	}
	ids = ids.AddAll(rIDs)
	return outputShape, action, ids, nil
}

func Compose(pos lexer.Position, l Expression, r Expression) Expression {
	if l == nil {
		return r
	}
	return &CompositionExpression{pos, l, r}
}
