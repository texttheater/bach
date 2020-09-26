package functions

import (
	"bytes"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/parameters"
	"github.com/texttheater/bach/states"
	"github.com/texttheater/bach/types"
)

type TemplateLiteralExpression struct {
	Pos    lexer.Position
	Pieces []Expression
}

func (x TemplateLiteralExpression) Position() lexer.Position {
	return x.Pos
}

func (x *TemplateLiteralExpression) Typecheck(inputShape Shape, params []*parameters.Parameter) (Shape, states.Action, *states.IDStack, error) {
	if len(params) > 0 {
		return Shape{}, nil, nil, errors.E(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	outputShape := Shape{
		Type:  types.StrType{},
		Stack: inputShape.Stack,
	}
	pieceActions := make([]states.Action, len(x.Pieces))
	for i, piece := range x.Pieces {
		_, pieceAction, _, err := piece.Typecheck(inputShape, nil)
		if err != nil {
			return Shape{}, nil, nil, err
		}
		pieceActions[i] = pieceAction
	}
	action := func(inputState states.State, args []states.Action) *states.Thunk {
		buffer := bytes.Buffer{}
		for _, pieceAction := range pieceActions {
			pieceThunk := pieceAction(inputState, nil)
			pieceResult := pieceThunk.Eval()
			if pieceResult.Error != nil {
				return states.ThunkFromError(pieceResult.Error)
			}
			out, err := pieceResult.Value.Out()
			if err != nil {
				return states.ThunkFromError(err)
			}
			buffer.WriteString(out)
		}
		return states.ThunkFromValue(states.StrValue(buffer.String()))
	}
	return outputShape, action, nil, nil
}
