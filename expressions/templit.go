package expressions

import (
	"bytes"

	"github.com/alecthomas/participle/lexer"
	"github.com/texttheater/bach/errors"
	"github.com/texttheater/bach/params"
	"github.com/texttheater/bach/shapes"
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

func (x *TemplateLiteralExpression) Typecheck(inputShape shapes.Shape, params []*params.Param) (shapes.Shape, states.Action, *states.IDStack, error) {
	if len(params) > 0 {
		return shapes.Shape{}, nil, nil, errors.TypeError(
			errors.Code(errors.ParamsNotAllowed),
			errors.Pos(x.Pos),
		)
	}
	outputShape := shapes.Shape{
		Type:  types.Str{},
		Stack: inputShape.Stack,
	}
	pieceActions := make([]states.Action, len(x.Pieces))
	for i, piece := range x.Pieces {
		_, pieceAction, _, err := piece.Typecheck(inputShape, nil)
		if err != nil {
			return shapes.Shape{}, nil, nil, err
		}
		pieceActions[i] = pieceAction
	}
	action := func(inputState states.State, args []states.Action) *states.Thunk {
		buffer := bytes.Buffer{}
		for _, pieceAction := range pieceActions {
			val, err := pieceAction(inputState, nil).Eval()
			if err != nil {
				return states.ThunkFromError(err)
			}
			out, err := val.Str()
			if err != nil {
				return states.ThunkFromError(err)
			}
			buffer.WriteString(out)
		}
		output := states.StrValue(buffer.String())
		return states.ThunkFromState(inputState.Replace(output))
	}
	return outputShape, action, nil, nil
}
