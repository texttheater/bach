package functions

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/parameters"
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

func (x CompositionExpression) Typecheck(inputShape Shape, params []*parameters.Parameter) (Shape, states.Action, *states.IDStack, error) {
	if len(params) > 0 {
		return Shape{}, nil, nil, states.E(
			states.Code(states.ParamsNotAllowed),
			states.Pos(x.Pos))

	}
	middleShape, lAction, ids, err := x.Left.Typecheck(inputShape, nil)
	if err != nil {
		return Shape{}, nil, nil, err
	}
	if (types.VoidType{}).Subsumes(middleShape.Type) {
		return Shape{}, nil, nil, states.E(
			states.Code(states.ComposeWithVoid),
			states.Pos(x.Right.Position()))

	}
	outputShape, rAction, rIDs, err := x.Right.Typecheck(middleShape, nil)
	if err != nil {
		return Shape{}, nil, nil, err
	}
	action := func(inputState states.State, args []states.Action) *states.Thunk {
		thunk := lAction(inputState, nil)
		res := thunk.Eval()
		if res.Error != nil {
			return states.ThunkFromError(res.Error)
		}
		if res.Drop {
			return &states.Thunk{
				Result: states.Result{
					Drop: true,
				},
			}
		}
		state := states.State{
			Value:     res.Value,
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
